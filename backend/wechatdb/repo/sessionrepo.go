package repo

import (
	"context"
	"github.com/labstack/gommon/log"
	"strings"
	"sync"
	"wechat-secretary/backend/wechatdb/model"
)
import datasource "wechat-secretary/backend/wechatdb/datasource/v4"

type SessionRepo struct {
	DataSource  *datasource.DataSource
	Mux         sync.RWMutex
	ChangeCh    chan struct{}
	ContactRepo *ContactRepo
}

func NewSessionRepo(path string, repo *ContactRepo) (*SessionRepo, error) {
	source, err := datasource.New(path)
	if err != nil {
		log.Errorf("Failed to init wechatdb datasource err:%s", err)
	}

	changeChan := source.GetDbChangeChan(datasource.Session)
	sr := &SessionRepo{
		DataSource:  source,
		ContactRepo: repo,
		ChangeCh:    changeChan,
	}

	go sr.watch()
	return sr, nil
}

func (sr *SessionRepo) GetSessions(context context.Context, key string, limit, offset int) []*model.SessionVo {
	sessions, err := sr.DataSource.GetSessions(context, key, limit, offset)
	if err != nil {
		log.Errorf("Failed to get session,err:%s", err)
	}

	svs := make([]*model.SessionVo, 0)
	for _, session := range sessions {
		wrap := session.Wrap()
		if strings.HasSuffix(wrap.UserName, "@chatroom") {
			chatRoom := sr.ContactRepo.GetChatRoomByUsername(context, wrap.UserName)
			if chatRoom != nil && len(chatRoom.NickName) != 0 {
				wrap.NickName = chatRoom.NickName
				wrap.HeadUrl = chatRoom.HeadUrl
			}
		} else {
			contact := sr.ContactRepo.GetContactModelByUsername(context, wrap.UserName)
			if contact != nil {
				wrap.NickName = contact.NickName
				wrap.HeadUrl = contact.SmallHeadUrl
			}
		}
		svs = append(svs, wrap)
	}
	return svs
}

func (sr *SessionRepo) watch() {
	for {
		select {
		case _, ok := <-sr.ChangeCh:
			if !ok {
				return
			} else {
				//
			}
		}
	}
}
