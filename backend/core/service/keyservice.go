package service

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
	"os"
	"sync"
	"time"
	"wechat-secretary/backend/core/storage"
	"wechat-secretary/backend/pkg/errorx"
	"wechat-secretary/backend/pkg/filemonitor"
	"wechat-secretary/backend/wechat"
	"wechat-secretary/backend/wechat/types"
)

type KeyService struct {
	weChatManager       *wechat.WechatManager
	weChatAccountKeyMap map[string]*types.WeChatAccountKey
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
		weChatManager:       weChatManager,
		lastEvents:          make(map[string]time.Time),
		pendingActions:      make(map[string]bool),
		weChatAccountKeyMap: make(map[string]*types.WeChatAccountKey),
		fs:                  client,
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
		s.weChatAccountKeyMap[key.Account] = key
	}

}

func (s *KeyService) GetKeyByAccount(name string) (*types.WeChatAccountKey, bool) {
	if r, exist := s.weChatAccountKeyMap[name]; exist {
		return r, true
	}
	return nil, false
}

func (s *KeyService) GetKeys() []*types.WeChatAccountKey {
	keys := make([]*types.WeChatAccountKey, 0)

	for _, key := range s.weChatAccountKeyMap {
		keys = append(keys, key)
	}

	return keys
}

func (s *KeyService) DecryptKey() error {
	accounts := s.weChatManager.GetAccounts()
	if len(accounts) == 0 {
		return errorx.WeChatProcessNotExist()
	}

	if len(accounts) == 1 {
		// key, imgKey := m.ctx.DataKey, m.ctx.ImgKey
		// if len(key) == 0 || len(imgKey) == 0 || force {
		// 	key, imgKey, err = m.ctx.WeChatInstances[0].GetKey(context.Background())
		// 	if err != nil {
		// 		return "", err
		// 	}
		// 	m.ctx.Refresh()
		// 	m.ctx.UpdateConfig()
		// }
		account, err := s.weChatManager.GetAccount(accounts[0].Name)
		if err != nil {
			log.Errorf("Failed to get WeChat process account")
			return errorx.WeChatProcessNotExist()
		}

		key, err := s.weChatManager.GetAccountKey(account.Name)

		if err != nil {
			log.Errorf("Failed to get account key")
		}

		_, err = s.fs.Create(key.Account, key)
		if err != nil {
			return nil
		}
		// result := fmt.Sprintf("Data Key: [%s]\nImage Key: [%s]", key, imgKey)
		// if m.ctx.Version == 4 && showXorKey {
		// 	if b, err := dat2img.ScanAndSetXorKey(m.ctx.DataDir); err == nil {
		// 		result += fmt.Sprintf("\nXor Key: [0x%X]", b)
		// 	}
		// }
		return nil
	}
	return nil
	// if pid == 0 {
	// 	str := "Select a process:\n"
	// 	for _, ins := range m.ctx.WeChatInstances {
	// 		str += fmt.Sprintf("PID: %d. %s[Version: %s Data Dir: %s ]\n", ins.PID, ins.Name, ins.FullVersion, ins.DataDir)
	// 	}
	// 	return str, nil
	// }
	// for _, ins := range m.ctx.WeChatInstances {
	// 	if ins.PID == uint32(pid) {
	// 		key, imgKey := ins.Key, ins.ImgKey
	// 		if len(key) == 0 || len(imgKey) == 0 || force {
	// 			key, imgKey, err = ins.GetKey(context.Background())
	// 			if err != nil {
	// 				return "", err
	// 			}
	// 			m.ctx.Refresh()
	// 			m.ctx.UpdateConfig()
	// 		}
	// 		result := fmt.Sprintf("Data Key: [%s]\nImage Key: [%s]", key, imgKey)
	// 		if m.ctx.Version == 4 && showXorKey {
	// 			if b, err := dat2img.ScanAndSetXorKey(m.ctx.DataDir); err == nil {
	// 				result += fmt.Sprintf("\nXor Key: [0x%X]", b)
	// 			}
	// 		}
	// 		return result, nil
	// 	}
	// }

}

func (s *KeyService) loadWechatAccountKey() ([]*types.WeChatAccountKey, error) {
	objs, err := s.fs.List(storage.WechatAccountKey)
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
