package service

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"sync"
	"time"
	"wechat-secretary/backend/pkg/errorx"
	"wechat-secretary/backend/pkg/filemonitor"
	"wechat-secretary/backend/pkg/logs"
	"wechat-secretary/backend/pkg/util"
	"wechat-secretary/backend/wechat/decrypt"
	"wechat-secretary/backend/wechat/types"
)

var (
	DebounceTime = 1 * time.Second
	MaxWaitTime  = 10 * time.Second
)

type Config struct {
	DataKey  string    `json:"dataKey"`
	DataDir  string    `json:"dataDir"`
	WorkDir  string    `json:"workDir"`
	Platform string    `json:"platform"`
	Version  string    `json:"version"`
	Time     time.Time `json:"time"`
}

type CoreService struct {
	Config         *Config
	KeyService     *KeyService
	mutex          sync.Mutex
	wechatAccounts map[string]*WechatAccountManager
}

func NewCoreService(service *KeyService) *CoreService {
	return &CoreService{
		Config:         nil,
		KeyService:     service,
		wechatAccounts: map[string]*WechatAccountManager{},
	}
}

func (s *CoreService) Decrypt(accountName string) error {
	account, b := s.KeyService.GetAccountDetailsByAccountName(accountName)
	if !b {
		logs.Errorf("Failed to get wechat account")
		return nil
	}

	wam := &WechatAccountManager{
		weChatAccount:  account,
		lastEvents:     map[string]time.Time{},
		pendingActions: map[string]bool{},
	}

	err := wam.Decrypt()
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.wechatAccounts[accountName] = wam

	return nil
}

type WechatAccountManager struct {
	weChatAccount  *types.WeChatAccountDetails
	lastEvents     map[string]time.Time
	pendingActions map[string]bool
	fm             *filemonitor.FileMonitor
	mutex          sync.Mutex
}

func (w *WechatAccountManager) Decrypt() error {
	dbGroup, err := filemonitor.NewFileGroup("wechat", w.weChatAccount.DataSourceDir, `.*\.db$`, []string{"fts"})
	if err != nil {
		return err
	}

	dbFiles, err := dbGroup.List()
	if err != nil {
		return err
	}

	for _, dbFile := range dbFiles {
		if err := w.DecryptDBFile(dbFile); err != nil {
			logs.Errorf("DecryptDBFile %s failed: %v", dbFile, err)
			continue
		}
	}

	// 开启自动解密
	err = w.StartAutoDecrypt()
	if err != nil {
		logs.Errorf("Failed to auto decrypt")
		return err
	}
	return nil
}

func (w *WechatAccountManager) DecryptDBFile(dbFile string) error {

	decryptor, err := decrypt.NewDecryptor(w.weChatAccount.Platform, w.weChatAccount.Version)
	if err != nil {
		return err
	}

	output := filepath.Join(w.weChatAccount.DataDecryptDir, dbFile[len(w.weChatAccount.DataSourceDir):])
	if err := util.PrepareDir(filepath.Dir(output)); err != nil {
		return err
	}

	outputTemp := output + ".tmp"
	outputFile, err := os.Create(outputTemp)
	if err != nil {
		logs.Errorf("Failed to create output file: %v", err)
	}
	defer func() {
		outputFile.Close()
		if err := os.Rename(outputTemp, output); err != nil {
			logs.Errorf("Failed to rename %s to %s", outputTemp, output)
		}
	}()

	if err := decryptor.Decrypt(context.Background(), dbFile, w.weChatAccount.DataKey, outputFile); err != nil {
		if err == errorx.ErrAlreadyDecrypted() {
			if data, err := os.ReadFile(dbFile); err == nil {
				outputFile.Write(data)
			}
			return nil
		}
		logs.Errorf("Failed to decrypt %s", dbFile)
		return err
	}

	logs.Infof("Decrypted %s to %s", dbFile, output)
	return nil
}

func (w *WechatAccountManager) StartAutoDecrypt() error {
	logs.Infof("Start auto decrypt, data dir: %s", w.weChatAccount.DataSourceDir)
	dbGroup, err := filemonitor.NewFileGroup("wechat", w.weChatAccount.DataSourceDir, `.*\.db$`, []string{"fts"})
	if err != nil {
		return err
	}
	dbGroup.AddCallback(w.DecryptFileCallback)

	w.fm = filemonitor.NewFileMonitor()
	w.fm.AddGroup(dbGroup)
	if err := w.fm.Start(); err != nil {
		logs.Errorf("Failed to start file monitor, data dir: %s", w.weChatAccount.DataSourceDir)
		return err
	}
	return nil
}

func (w *WechatAccountManager) DecryptFileCallback(event fsnotify.Event) error {
	// Local file system
	// WRITE         "/db_storage/message/message_0.db"
	// WRITE         "/db_storage/message/message_0.db"
	// WRITE|CHMOD   "/db_storage/message/message_0.db"
	// Syncthing
	// REMOVE        "/app/data/db_storage/session/session.db"
	// CREATE        "/app/data/db_storage/session/session.db" ← "/app/data/db_storage/session/.syncthing.session.db.tmp"
	// CHMOD         "/app/data/db_storage/session/session.db"
	if !(event.Op.Has(fsnotify.Write) || event.Op.Has(fsnotify.Create)) {
		return nil
	}

	w.mutex.Lock()
	w.lastEvents[event.Name] = time.Now()

	if !w.pendingActions[event.Name] {
		w.pendingActions[event.Name] = true
		w.mutex.Unlock()
		go w.waitAndProcess(event.Name)
	} else {
		w.mutex.Unlock()
	}

	return nil
}

func (w *WechatAccountManager) waitAndProcess(dbFile string) {
	start := time.Now()
	for {
		time.Sleep(DebounceTime)

		w.mutex.Lock()
		lastEventTime := w.lastEvents[dbFile]
		elapsed := time.Since(lastEventTime)
		totalElapsed := time.Since(start)

		if elapsed >= DebounceTime || totalElapsed >= MaxWaitTime {
			w.pendingActions[dbFile] = false
			w.mutex.Unlock()
			logs.Infof("Processing file: %s", dbFile)
			w.DecryptDBFile(dbFile)
			return
		}
		w.mutex.Unlock()
	}
}
