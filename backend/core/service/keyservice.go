package service

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
	"os"
	"path/filepath"
	"sync"
	"time"
	"wechat-secretary/backend/core/storage"
	"wechat-secretary/backend/pkg/errorx"
	"wechat-secretary/backend/pkg/filemonitor"
	"wechat-secretary/backend/wechat"
	"wechat-secretary/backend/wechat/types"
)

type KeyService struct {
	weChatManager    *wechat.WechatManager
	weChatAccountMap map[string]*types.WeChatAccount
	// conf           Config
	lastEvents     map[string]time.Time
	pendingActions map[string]bool
	mutex          sync.Mutex
	fm             *filemonitor.FileMonitor
	fs             *storage.FsClient
}

func NewKeyService(weChatManager *wechat.WechatManager) *KeyService {
	client := &storage.FsClient{}
	client.Init(storage.StoreGroupWechatAccountKey)
	return &KeyService{
		weChatManager:    weChatManager,
		lastEvents:       make(map[string]time.Time),
		pendingActions:   make(map[string]bool),
		weChatAccountMap: make(map[string]*types.WeChatAccount),
		fs:               client,
	}
}

func (s *KeyService) Init() {
	keys, err := s.loadWechatAccountKey()
	if err != nil {
		log.Errorf("Failed to get wechat account key")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, key := range keys {
		account, err := s.weChatManager.GetAccount(key.Account)
		if err != nil {
			log.Errorf("Failed to get wechat account process")
		} else {
			wa := &types.WeChatAccount{
				WeChatAccountKey: key,
				Platform:         account.Platform,
				Version:          account.Version,
				FullVersion:      account.FullVersion,
				DataDir:          account.DataDir,
				ExePath:          account.ExePath,
				Status:           account.Status,
			}
			s.weChatAccountMap[key.Account] = wa
		}
	}

}

func (s *KeyService) GetKeyByAccount(name string) (*types.WeChatAccount, bool) {
	if r, exist := s.weChatAccountMap[name]; exist {
		return r, true
	}
	return nil, false
}

func (s *KeyService) GetKeys() []*types.WeChatAccountKey {
	keys := make([]*types.WeChatAccountKey, 0)

	for _, key := range s.weChatAccountMap {
		keys = append(keys, key.WeChatAccountKey)
	}

	return keys
}

func (s *KeyService) GetAccounts() []*types.WeChatAccount {
	keys := make([]*types.WeChatAccount, 0)

	for _, key := range s.weChatAccountMap {
		keys = append(keys, key)
	}

	return keys
}

func (s *KeyService) DecryptKey() error {
	accounts := s.weChatManager.GetAccounts()
	if len(accounts) == 0 {
		return errorx.WeChatProcessNotExist()
	}

	for i, _ := range accounts {
		account, err := s.weChatManager.GetAccount(accounts[i].Name)
		if err != nil {
			log.Errorf("Failed to get WeChat process account")
			return errorx.WeChatProcessNotExist()
		}

		key, err := s.weChatManager.GetAccountKey(account.Name)

		if err != nil {
			log.Errorf("Failed to get account key")
			return err
		}

		// 默认设置数据目录
		key.DataDir = filepath.Join(storage.WechatSecretaryPrefix, key.Account)

		_, err = s.fs.Create(key.Account, key)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (s *KeyService) loadWechatAccountKey() ([]*types.WeChatAccountKey, error) {
	objs, err := s.fs.List("")
	if err != nil {
		return nil, err
	}

	var ret []*types.WeChatAccountKey
	if files, ok := objs.([]*storage.FileInfo); ok {
		for _, file := range files {
			func() {
				obj := types.WeChatAccountKey{}
				f, err := os.Open(file.Path)
				defer f.Close()
				if err != nil {
					log.Errorf("Failed to open", "file", file.Path, "err", err)
					return
				}
				if err = json.NewDecoder(f).Decode(&obj); err != nil {
					log.Errorf("Failed to unmarshal", "file", file.Path, "err", err)
					return
				}
				ret = append(ret, &obj)
			}()
		}
	}
	return ret, nil
}
