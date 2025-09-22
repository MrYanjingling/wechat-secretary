package wechat

import (
	"context"
	"github.com/labstack/gommon/log"
	"runtime"
	"wechat-secretary/backend/pkg/errorx"
	"wechat-secretary/backend/wechat/decrypt"
	"wechat-secretary/backend/wechat/key"
	"wechat-secretary/backend/wechat/model"
	"wechat-secretary/backend/wechat/process"
	"wechat-secretary/backend/wechat/types"
)

// Manager 微信管理器
type WechatManager struct {
	detector   process.Detector
	accounts   []*Account
	processMap map[string]*model.Process
}

// Account 表示一个微信账号
type Account struct {
	Name        string
	Platform    string
	Version     int
	FullVersion string
	DataDir     string
	Key         string
	ImgKey      string
	PID         uint32
	ExePath     string
	Status      string
}

// RefreshStatus 刷新账号的进程状态
// func (a *Account) RefreshStatus() error {
// 	// 查找所有微信进程
// 	Load()
//
// 	process, err := GetProcess(a.Name)
// 	if err != nil {
// 		a.Status = model.StatusOffline
// 		return nil
// 	}
//
// 	if process.AccountName == a.Name {
// 		// 更新进程信息
// 		a.PID = process.PID
// 		a.ExePath = process.ExePath
// 		a.Platform = process.Platform
// 		a.Version = process.Version
// 		a.FullVersion = process.FullVersion
// 		a.Status = process.Status
// 		a.DataDir = process.DataDir
// 	}
//
// 	return nil
// }

// GetKey 获取账号的密钥
// func (a *Account) GetKey(ctx context.Context) (string, string, error) {
// 	// 如果已经有密钥，直接返回
// 	if a.Key != "" && (a.ImgKey != "" || a.Version == 3) {
// 		return a.Key, a.ImgKey, nil
// 	}
//
// 	// 刷新进程状态
// 	if err := a.RefreshStatus(); err != nil {
// 		return "", "", errors.RefreshProcessStatusFailed(err)
// 	}
//
// 	// 检查账号状态
// 	if a.Status != model.StatusOnline {
// 		return "", "", errors.WeChatAccountNotOnline(a.Name)
// 	}
//
// 	// 创建密钥提取器 - 使用新的接口，传入平台和版本信息
// 	extractor, err := key.NewExtractor(a.Platform, a.Version)
// 	if err != nil {
// 		return "", "", err
// 	}
//
// 	process, err := GetProcess(a.Name)
// 	if err != nil {
// 		return "", "", err
// 	}
//
// 	validator, err := decrypt.NewValidator(process.Platform, process.Version, process.DataDir)
// 	if err != nil {
// 		return "", "", err
// 	}
//
// 	extractor.SetValidate(validator)
//
// 	// 提取密钥
// 	dataKey, imgKey, err := extractor.Extract(ctx, process)
// 	if err != nil {
// 		return "", "", err
// 	}
//
// 	if dataKey != "" {
// 		a.Key = dataKey
// 	}
//
// 	if imgKey != "" {
// 		a.ImgKey = imgKey
// 	}
//
// 	return dataKey, imgKey, nil
// }

// DecryptDatabase 解密数据库
// func (a *Account) DecryptDatabase(ctx context.Context, dbPath, outputPath string) error {
// 	// 获取密钥
// 	hexKey, _, err := a.GetKey(ctx)
// 	if err != nil {
// 		return err
// 	}
//
// 	// 创建解密器 - 传入平台信息和版本
// 	decryptor, err := decrypt.NewDecryptor(a.Platform, a.Version)
// 	if err != nil {
// 		return err
// 	}
//
// 	// 创建输出文件
// 	output, err := os.Create(outputPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer output.Close()
//
// 	// 解密数据库
// 	return decryptor.Decrypt(ctx, dbPath, hexKey, output)
// }

// NewManager 创建新的微信管理器
func NewManager() *WechatManager {
	return &WechatManager{
		detector:   process.NewDetector(runtime.GOOS),
		accounts:   make([]*Account, 0),
		processMap: make(map[string]*model.Process),
	}
}

// Init 加载微信进程信息
func (m *WechatManager) Init() error {
	// 查找微信进程
	processes, err := m.detector.FindProcesses()
	if err != nil {
		return err
	}

	// 转换为账号信息
	accounts := make([]*Account, 0, len(processes))
	processMap := make(map[string]*model.Process, len(processes))

	for _, p := range processes {
		account := NewAccount(p)

		accounts = append(accounts, account)
		if account.Name != "" {
			processMap[account.Name] = p
		}
	}

	m.accounts = accounts
	m.processMap = processMap

	return nil
}

// GetAccount 获取指定名称的账号
func (m *WechatManager) GetAccount(name string) (*Account, error) {
	p, err := m.GetProcess(name)
	if err != nil {
		return nil, err
	}
	return NewAccount(p), nil
}

func (m *WechatManager) GetProcess(name string) (*model.Process, error) {
	p, ok := m.processMap[name]
	if !ok {
		return nil, errorx.WeChatAccountNotFound(name)
	}
	return p, nil
}

// GetAccounts 获取所有账号
func (m *WechatManager) GetAccounts() []*Account {
	return m.accounts
}

// DecryptDatabase 便捷方法：通过账号名解密数据库
// func (m *WechatManager) DecryptDatabase(ctx context.Context, accountName, dbPath, outputPath string) error {
// 	// 获取账号
// 	account, err := m.GetAccount(accountName)
// 	if err != nil {
// 		return err
// 	}
//
// 	// 使用账号解密数据库
// 	return account.DecryptDatabase(ctx, dbPath, outputPath)
// }

func (m *WechatManager) GetAccountKey(name string) (*types.WeChatAccountDetails, error) {

	// 刷新进程状态
	if err := m.RefreshStatus(name); err != nil {
		return nil, errorx.WeChatProcessNotExist()
	}

	process, _ := m.GetProcess(name)

	// 创建密钥提取器 - 使用新的接口，传入平台和版本信息
	extractor, err := key.NewExtractor(process.Platform, process.Version)
	if err != nil {
		return nil, err
	}

	validator, err := decrypt.NewValidator(process.Platform, process.Version, process.DataDir)
	if err != nil {
		return nil, err
	}

	extractor.SetValidate(validator)

	// 提取密钥
	dataKey, imgKey, err := extractor.Extract(context.Background(), process)
	if err != nil {
		return nil, err
	}

	return &types.WeChatAccountDetails{
		Account:       name,
		DataKey:       dataKey,
		ImgKey:        imgKey,
		Platform:      process.Platform,
		Version:       process.Version,
		FullVersion:   process.FullVersion,
		DataSourceDir: process.DataDir,
		ExePath:       process.ExePath,
		NonLocal:      false,
		Status:        process.Status,
	}, nil

}

func (w *WechatManager) RefreshStatus(accountName string) error {

	process, err := w.GetProcess(accountName)
	if err != nil {
		process.Status = model.StatusOffline
		log.Errorf("Failed to find wechat process")
	}

	// if process.AccountName == accountName {
	// 	// 更新进程信息
	// 	a.PID = process.PID
	// 	a.ExePath = process.ExePath
	// 	a.Platform = process.Platform
	// 	a.Version = process.Version
	// 	a.FullVersion = process.FullVersion
	// 	a.Status = process.Status
	// 	a.DataDir = process.DataDir
	// }

	return nil
}

// NewAccount 创建新的账号对象
func NewAccount(proc *model.Process) *Account {
	return &Account{
		Name:        proc.AccountName,
		Platform:    proc.Platform,
		Version:     proc.Version,
		FullVersion: proc.FullVersion,
		DataDir:     proc.DataDir,
		PID:         proc.PID,
		ExePath:     proc.ExePath,
		Status:      proc.Status,
	}
}
