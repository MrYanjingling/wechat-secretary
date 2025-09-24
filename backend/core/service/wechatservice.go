package service

import (
	"sync"
	"wechat-secretary/backend/wechat/types"
	"wechat-secretary/backend/wechatdb/repo"
)

type WechatService struct {
	account     *types.WeChatAccountDetails
	MediaRepo   *repo.MediaRepo
	ContactRepo *repo.ContactRepo
	SessionRepo *repo.SessionRepo
	mutex       sync.Mutex
}

func NewWechatService(account *types.WeChatAccountDetails) (*WechatService, error) {
	contactRepo, err := repo.NewContactRepo(account.DataDecryptDir)
	if err != nil {
		return nil, err
	}

	sessionRepo, err := repo.NewSessionRepo(account.DataDecryptDir, contactRepo)
	if err != nil {
		return nil, err
	}

	mediaRepo, err := repo.NewMediaRepo(account.DataDecryptDir, contactRepo)
	if err != nil {
		return nil, err
	}

	return &WechatService{
		account:     account,
		ContactRepo: contactRepo,
		MediaRepo:   mediaRepo,
		SessionRepo: sessionRepo,
	}, nil
}
