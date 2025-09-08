package repo

import (
	"context"
	"github.com/labstack/gommon/log"
	"strings"
	"sync"
	"wechat-secretary/backend/wechat/model"
	v4 "wechat-secretary/backend/wechat/model/v4"
)
import datasource "wechat-secretary/backend/wechat/datasource/v4"

type ContactRepo struct {
	DataSource      *datasource.DataSource
	Mux             sync.RWMutex
	ChangeCh        chan struct{}
	ContactLabelMap map[string]*v4.ContactLabelModel

	ContactModelMap map[string]*v4.ContactModel // map[username]
	ContactMap      map[string]*model.ContactVo
	AliasMap        map[string][]*model.ContactVo
	RemarkMap       map[string][]*model.ContactVo
	NickNameMap     map[string][]*model.ContactVo

	ChatRoomModelMap    map[string]*v4.ChatRoomModel // map[username]
	ChatRoomMap         map[string]*model.ChatRoom
	RemarkChatRoomMap   map[string][]*model.ChatRoom
	NickNameChatRoomMap map[string][]*model.ChatRoom
}

func NewContactRepo(path string) (*ContactRepo, error) {
	source, err := datasource.New(path)
	if err != nil {
		log.Errorf("Failed to init wechat datasource err:%s", err)
	}

	changeChan := source.GetDbChangeChan(datasource.Contact)
	cr := &ContactRepo{
		DataSource:      source,
		ContactLabelMap: make(map[string]*v4.ContactLabelModel, 0),
		ContactModelMap: make(map[string]*v4.ContactModel, 0),
		ContactMap:      make(map[string]*model.ContactVo),
		AliasMap:        make(map[string][]*model.ContactVo),
		RemarkMap:       make(map[string][]*model.ContactVo),
		NickNameMap:     make(map[string][]*model.ContactVo),
		ChangeCh:        changeChan,
	}

	cr.Index()

	go cr.watch()
	return cr, nil
}

func (ct *ContactRepo) Index() {
	ct.indexContactLabels()
	ct.indexContacts()
	ct.indexChatRooms()
}

func (ct *ContactRepo) indexChatRooms() {
	ct.Mux.RLock()
	defer ct.Mux.RUnlock()

	chatRoomModelMap := make(map[string]*v4.ChatRoomModel, 0)
	chatRoomMap := make(map[string]*model.ChatRoom, 0)
	remarkChatRoomMap := make(map[string][]*model.ChatRoom, 0)
	nickNameChatRoomMap := make(map[string][]*model.ChatRoom, 0)

	rooms, err := ct.DataSource.GetAllChatRooms(context.Background())
	if err != nil {
		log.Errorf("Failed to get chatroom err:%s", err)
	}

	for _, chatRoomModel := range rooms {
		chatRoomModelMap[chatRoomModel.UserName] = chatRoomModel

		cm := chatRoomModel.Wrap()
		contactModel := ct.ContactModelMap[cm.Name]
		cm.NickName = contactModel.NickName
		cm.Remark = contactModel.Remark
		chatRoomMap[cm.Name] = cm

		if cm.Remark != "" {
			remark, ok := remarkChatRoomMap[cm.Remark]
			if !ok {
				remark = make([]*model.ChatRoom, 0)
			}
			remark = append(remark, cm)
			remarkChatRoomMap[cm.Remark] = remark
		}
		if cm.NickName != "" {
			nickName, ok := nickNameChatRoomMap[cm.NickName]
			if !ok {
				nickName = make([]*model.ChatRoom, 0)
			}
			nickName = append(nickName, cm)
			nickNameChatRoomMap[cm.NickName] = nickName
		}
	}

	ct.ChatRoomModelMap = chatRoomModelMap
	ct.ChatRoomMap = chatRoomMap
	ct.RemarkChatRoomMap = remarkChatRoomMap
	ct.NickNameChatRoomMap = nickNameChatRoomMap

}

func (ct *ContactRepo) indexContactLabels() {
	contactLabelMap := make(map[string]*v4.ContactLabelModel, 0)
	ct.Mux.RLock()
	defer ct.Mux.RUnlock()
	label, err := ct.DataSource.GetContactLabel(context.Background())
	if err != nil {
		log.Errorf("Failed to get contact label err:%s", err)
	}
	for _, model := range label {
		contactLabelMap[model.LabelId] = model
	}
	// 赋值
	ct.ContactLabelMap = contactLabelMap
}

func (ct *ContactRepo) indexContacts() {
	ct.Mux.RLock()
	defer ct.Mux.RUnlock()

	contactModelMap := make(map[string]*v4.ContactModel, 0)
	aliasMap := make(map[string][]*model.ContactVo, 0)
	contactMap := make(map[string]*model.ContactVo)
	remarkMap := make(map[string][]*model.ContactVo)
	nickNameMap := make(map[string][]*model.ContactVo)

	contacts, err := ct.DataSource.GetAllContact(context.Background())
	if err != nil {
		log.Errorf("Failed to get contacts err:%s", err)
	}

	for _, contactModel := range contacts {
		contactModelMap[contactModel.UserName] = contactModel
		if strings.HasSuffix(contactModel.UserName, "@chatroom") {
			continue
		}
		if contactModel.LocalType == 1 {
			// 朋友
			ct.indexContact(aliasMap, contactMap, remarkMap, nickNameMap, contactModel)
			continue
		} else if contactModel.LocalType == 5 {
			// 企业微信
			ct.indexContactWework(aliasMap, contactMap, remarkMap, nickNameMap, contactModel)
			continue
		} else {
			continue
		}
	}

	ct.ContactModelMap = contactModelMap
	ct.ContactMap = contactMap
	ct.AliasMap = aliasMap
	ct.NickNameMap = nickNameMap
	ct.RemarkMap = remarkMap
}

func (ct *ContactRepo) indexContact(aliasMap map[string][]*model.ContactVo, contactMap map[string]*model.ContactVo, remarkMap map[string][]*model.ContactVo, nickNameMap map[string][]*model.ContactVo, contactModel *v4.ContactModel) {
	vo := contactModel.WrapVo()
	contactMap[vo.UserName] = vo

	if vo.Alias != "" {
		alias, ok := aliasMap[vo.Alias]
		if !ok {
			alias = make([]*model.ContactVo, 0)
		}
		alias = append(alias, vo)
		aliasMap[vo.Alias] = alias
	}
	if vo.Remark != "" {
		remark, ok := remarkMap[vo.Remark]
		if !ok {
			remark = make([]*model.ContactVo, 0)
		}
		remark = append(remark, vo)
		remarkMap[vo.Remark] = remark
	}
	if vo.NickName != "" {
		nickName, ok := nickNameMap[vo.NickName]
		if !ok {
			nickName = make([]*model.ContactVo, 0)
		}
		nickName = append(nickName, vo)
		nickNameMap[vo.NickName] = nickName
	}

}

func (ct *ContactRepo) indexContactWework(aliasMap map[string][]*model.ContactVo, contactMap map[string]*model.ContactVo, remarkMap map[string][]*model.ContactVo, nickNameMap map[string][]*model.ContactVo, contactModel *v4.ContactModel) {
	vo := contactModel.WrapVo()
	contactMap[vo.UserName] = vo

	if vo.Alias != "" {
		alias, ok := aliasMap[vo.Alias]
		if !ok {
			alias = make([]*model.ContactVo, 0)
		}
		alias = append(alias, vo)
		aliasMap[vo.Alias] = alias
	}
	if vo.Remark != "" {
		remark, ok := remarkMap[vo.Remark]
		if !ok {
			remark = make([]*model.ContactVo, 0)
		}
		remark = append(remark, vo)
		remarkMap[vo.Remark] = remark
	}
	if vo.NickName != "" {
		nickName, ok := nickNameMap[vo.NickName]
		if !ok {
			nickName = make([]*model.ContactVo, 0)
		}
		nickName = append(nickName, vo)
		nickNameMap[vo.NickName] = nickName
	}
}

func (ct *ContactRepo) GetContactByUsername(username string) *model.ContactVo {
	if vo, ok := ct.ContactMap[username]; ok {
		return vo
	} else {
		return nil
	}
}

func (ct *ContactRepo) GetContactByNickName(nickname string) *model.ContactVo {
	if vo, ok := ct.NickNameMap[nickname]; ok {
		return vo[0]
	} else {
		return nil
	}
}

func (ct *ContactRepo) watch() {
	for {
		select {
		case _, ok := <-ct.ChangeCh:
			if !ok {
				return
			} else {
				ct.Index()
			}
		}
	}
}
