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
	weChatAccountMap map[string]*types.WeChatAccountDetails
	// conf           Config
	lastEvents     map[string]time.Time
	pendingActions map[string]bool
	fm             *filemonitor.FileMonitor
	mutex          sync.Mutex
	fs             *storage.FsClient
}

func NewKeyService(weChatManager *wechat.WechatManager) *KeyService {
	client := &storage.FsClient{}
	client.Init(storage.StoreGroupWechatAccountKey)
	return &KeyService{
		weChatManager:    weChatManager,
		lastEvents:       make(map[string]time.Time),
		pendingActions:   make(map[string]bool),
		weChatAccountMap: make(map[string]*types.WeChatAccountDetails),
		fs:               client,
	}
}

func (s *KeyService) Init() {
	keys, err := s.loadWechatAccountDetails()
	if err != nil {
		log.Errorf("Failed to get wechat account key")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, key := range keys {
		if a, err := s.weChatManager.GetAccount(key.Account); err == nil {
			key.DataSourceDir = a.DataDir
			key.Platform = a.Platform
			key.Version = a.Version
			key.FullVersion = a.FullVersion
			key.ExePath = a.ExePath
			key.Status = a.Status
		}
		s.weChatAccountMap[key.Account] = key
	}

}

func (s *KeyService) GetAccountDetailsByAccountName(name string) (*types.WeChatAccountDetails, bool) {
	if r, exist := s.weChatAccountMap[name]; exist {
		return r, true
	}
	return nil, false
}

func (s *KeyService) GetAccountDetails() []*types.WeChatAccountDetails {
	keys := make([]*types.WeChatAccountDetails, 0)

	for _, key := range s.weChatAccountMap {
		keys = append(keys, key)
	}

	return keys
}

func (s *KeyService) GetWxAccount() []*wechat.Account {
	return s.weChatManager.GetAccounts()
}

func (s *KeyService) GetWxAccountByName(name string) *wechat.Account {
	a, err := s.weChatManager.GetAccount(name)
	if err != nil {
		log.Errorf("Failed to get wechat account by name")
		return nil
	}
	return a
}

func (s *KeyService) DecryptKey(name string) error {
	account := s.GetWxAccountByName(name)
	if account == nil {
		return errorx.WeChatProcessNotExist()
	}

	ad, err := s.weChatManager.GetAccountKey(account.Name)

	if err != nil {
		log.Errorf("Failed to get account key")
		return err
	}

	// 默认设置数据目录
	ad.DataDecryptDir = filepath.Join(storage.WechatSecretaryPrefix, ad.Account)

	return s.StoreAccount(ad)
}

func (s *KeyService) StoreAccount(account *types.WeChatAccountDetails) error {
	details, exist := s.weChatAccountMap[account.Account]
	if !exist {
		s.weChatAccountMap[account.Account] = details
	}

	_, err := s.fs.Create(account.Account, account)
	if err != nil {
		return err
	}
	return nil
}

func (s *KeyService) loadWechatAccountDetails() ([]*types.WeChatAccountDetails, error) {
	objs, err := s.fs.List("")
	if err != nil {
		return nil, err
	}

	var ret []*types.WeChatAccountDetails
	if files, ok := objs.([]*storage.FileInfo); ok {
		for _, file := range files {
			func() {
				obj := types.WeChatAccountDetails{}
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
