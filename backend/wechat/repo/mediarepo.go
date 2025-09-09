package repo

import (
	"context"
	"github.com/labstack/gommon/log"
	"sync"
	"wechat-secretary/backend/wechat/model"
)
import datasource "wechat-secretary/backend/wechat/datasource/v4"

type MediaRepo struct {
	DataSource  *datasource.DataSource
	Mux         sync.RWMutex
	ChangeCh    chan struct{}
	ContactRepo *ContactRepo
}

func NewMediaRepo(path string, repo *ContactRepo) (*MediaRepo, error) {
	source, err := datasource.New(path)
	if err != nil {
		log.Errorf("Failed to init wechat datasource err:%s", err)
	}

	changeChan := source.GetDbChangeChan(datasource.Media)
	mr := &MediaRepo{
		DataSource:  source,
		ContactRepo: repo,
		ChangeCh:    changeChan,
	}

	go mr.watch()
	return mr, nil
}

func (mr *MediaRepo) GetMedia(ctx context.Context, _type string, key string) *model.MediaVo {
	media, err := mr.DataSource.GetMedia(ctx, _type, key)
	if err != nil {
		log.Errorf("Failed to get media err:%s", err)
	}
	wrap := media.Wrap()
	return wrap
}

func (mr *MediaRepo) GetVoice(ctx context.Context, key string) *model.VoiceVo {
	voiceData, err := mr.DataSource.GetVoice(ctx, key)
	if err != nil {
		log.Errorf("Failed to get media err:%s", err)
	}
	if voiceData != nil {
		return &model.VoiceVo{
			Key:  key,
			Data: voiceData,
		}
	}
	return nil
}

func (mr *MediaRepo) watch() {
	for {
		select {
		case _, ok := <-mr.ChangeCh:
			if !ok {
				return
			} else {
				//
			}
		}
	}
}
