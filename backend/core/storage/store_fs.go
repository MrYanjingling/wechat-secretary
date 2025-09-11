package storage

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/labstack/gommon/log"
	"os"
	"path/filepath"
	"sync"
	"wechat-secretary/backend/pkg/errorx"
)

type FsClient struct {
	storePath string
}

var (
	_ Storage = (*FsClient)(nil)

	doOnce sync.Once
)

func (fc *FsClient) Init(sg StoreGroup) {
	_, err := os.Stat(storePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(storePath, 0777); err != nil {
				log.Errorf("%s: %v", storePath, err)
			}
		} else {
			log.Errorf("%s: %v", storePath, err)
		}
	}

	fc.storePath = filepath.Join(storePath, StoreGroupToString[sg])

	_, err = os.Stat(fc.storePath)
	if os.IsNotExist(err) {
		absPath, _ := filepath.Abs(fc.storePath)
		log.Infof("Created", "path", absPath)
		if err = os.MkdirAll(fc.storePath, 0711); err != nil {
			log.Error(err)
		}
	} else if err != nil {
		log.Error(err)
	}

	doOnce.Do(func() {
		gob.Register(map[string]interface{}{})
		gob.Register([]interface{}{})
	})
}

func (fc *FsClient) Create(key string, obj interface{}) (interface{}, error) {
	f, err := os.OpenFile(filepath.Join(fc.storePath, key), os.O_CREATE|os.O_RDWR|os.O_EXCL, 0640)
	if err != nil {
		log.Errorf("Failed to create file", "err", err)
		return nil, err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(obj)
	if err != nil {
		log.Errorf("Failed to encode", "err", err)
		return nil, err
	}
	return obj, nil
}

func (fc *FsClient) Get(key string) (interface{}, error) {
	data, err := os.ReadFile(filepath.Join(fc.storePath, key))
	if err != nil {
		log.Errorf("Failed to read", "err", err)
		return nil, err
	}
	return data, nil
}

func (fc *FsClient) List(key string) (interface{}, error) {
	var files []*FileInfo
	err := filepath.Walk(filepath.Join(fc.storePath, key), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, &FileInfo{
				Path: path,
			})
		}
		return nil
	})
	if err != nil {
		log.Errorf("Failed to list", "err", err)
	}
	return files, nil
}

func (fc *FsClient) Delete(key string) (interface{}, error) {

	f, err := os.OpenFile(filepath.Join(fc.storePath, key), os.O_RDONLY, 0640)
	if err != nil {
		if os.IsNotExist(err) {
			log.Errorf("Failed to open file", "err", err)
			return nil, os.ErrNotExist
		} else if isEphemeralError(err) {
			log.Errorf("Failed to open file", "err", err)
			return nil, errorx.WriteConflict()
		}
	}
	defer f.Close()

	err = os.Remove(filepath.Join(fc.storePath, key))
	if err != nil {
		log.Errorf("Failed to remove", "err", err)
		return nil, errorx.Internal()
	}
	return nil, nil
}

func (fc *FsClient) Update(key string, obj interface{}) (interface{}, error) {
	f, err := os.OpenFile(filepath.Join(fc.storePath, key), os.O_RDWR, 0640)
	if err != nil {
		if os.IsNotExist(err) {
			log.Errorf("Failed to open file", "err", err)
			return nil, os.ErrNotExist
		} else if isEphemeralError(err) {
			log.Errorf("Failed to open file", "err", err)
			return nil, errorx.WriteConflict()
		}
	}
	defer f.Close()

	if err = f.Truncate(0); err != nil {
		log.Errorf("Failed to truncate", "err", err)
		return nil, errorx.Internal()
	}
	if _, err = f.Seek(0, 0); err != nil {
		log.Errorf("Failed to seek", "err", err)
		return nil, errorx.Internal()
	}
	err = json.NewEncoder(f).Encode(obj)
	if err != nil {
		log.Errorf("Failed to marshal", "err", err)
		return nil, errorx.Internal()
	}

	return obj, nil
}
