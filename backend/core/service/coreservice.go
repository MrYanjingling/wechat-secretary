package service

import (
	"context"
	"os"
	"path/filepath"
	"time"
	"wechat-secretary/backend/pkg/errorx"
	"wechat-secretary/backend/pkg/filemonitor"
	"wechat-secretary/backend/pkg/logs"
	"wechat-secretary/backend/pkg/util"
	"wechat-secretary/backend/wechat/decrypt"
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
	Config     *Config
	KeyService *KeyService
}

func NewCoreService(service *KeyService) *CoreService {
	return &CoreService{
		Config:     nil,
		KeyService: service,
	}
}

func (s *CoreService) Decrypt(accountName string) error {
	account, b := s.KeyService.GetAccountByAccountName(accountName)
	if !b {
		logs.Errorf("Failed to get wechat account")
		return nil
	}

	dbGroup, err := filemonitor.NewFileGroup("wechat", account.DataDir, `.*\.db$`, []string{"fts"})
	if err != nil {
		return err
	}

	dbFiles, err := dbGroup.List()
	if err != nil {
		return err
	}

	for _, dbFile := range dbFiles {
		if err := s.DecryptDBFile(dbFile, account.DataDir, account.WeChatAccountKey.DataDir, account.Platform, account.Version, account.WeChatAccountKey.DataKey); err != nil {
			logs.Errorf("DecryptDBFile %s failed: %v", dbFile, err)
			continue
		}
	}

	return nil
}

func (s *CoreService) DecryptDBFile(dbFile string, dataDir string, workDir string, platform string, version int, dataKey string) error {

	decryptor, err := decrypt.NewDecryptor(platform, version)
	if err != nil {
		return err
	}

	output := filepath.Join(workDir, dbFile[len(dataDir):])
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

	if err := decryptor.Decrypt(context.Background(), dbFile, dataKey, outputFile); err != nil {
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
