package v4

import (
	"google.golang.org/protobuf/proto"
	"wechat-secretary/backend/wechat/model"
	"wechat-secretary/backend/wechat/model/wxproto"
)

type ChatRoomModel struct {
	ID        int    `json:"id"`
	UserName  string `json:"username"`
	Owner     string `json:"owner"`
	ExtBuffer []byte `json:"ext_buffer"`
}

func (c *ChatRoomModel) Wrap() *model.ChatRoom {

	var users []model.ChatRoomUser
	if len(c.ExtBuffer) != 0 {
		users = c.ParseRoomData(c.ExtBuffer)
	}

	user2DisplayName := make(map[string]string, len(users))
	for _, user := range users {
		if user.DisplayName != "" {
			user2DisplayName[user.UserName] = user.DisplayName
		}
	}

	return &model.ChatRoom{
		Name:             c.UserName,
		Owner:            c.Owner,
		Users:            users,
		User2DisplayName: user2DisplayName,
	}
}

func (c *ChatRoomModel) ParseRoomData(b []byte) (users []model.ChatRoomUser) {
	var pbMsg wxproto.RoomData
	if err := proto.Unmarshal(b, &pbMsg); err != nil {
		return
	}
	if pbMsg.Users == nil {
		return
	}

	users = make([]model.ChatRoomUser, 0, len(pbMsg.Users))
	for _, user := range pbMsg.Users {
		u := model.ChatRoomUser{UserName: user.UserName}
		if user.DisplayName != nil {
			u.DisplayName = *user.DisplayName
		}
		users = append(users, u)
	}
	return users
}

type ContactLabelModel struct {
	LabelId   string `json:"label_id_"`
	LabelName string `json:"label_name_"`
	SortOrder string `json:"sort_order_"`
}

// CREATE TABLE contact(
// id INTEGER PRIMARY KEY,
// username TEXT,
// local_type INTEGER,
// alias TEXT,
// encrypt_username TEXT,
// flag INTEGER,
// delete_flag INTEGER,
// verify_flag INTEGER,
// remark TEXT,
// remark_quan_pin TEXT,
// remark_pin_yin_initial TEXT,
// nick_name TEXT,
// pin_yin_initial TEXT,
// quan_pin TEXT,
// big_head_url TEXT,
// small_head_url TEXT,
// head_img_md5 TEXT,
// chat_room_notify INTEGER,
// is_in_chat_room INTEGER,
// description TEXT,
// extra_buffer BLOB,
// chat_room_type INTEGER
// )
type ContactModel struct {
	UserName     string `json:"username"`   // 微信号id  gz_ 公众号
	Alias        string `json:"alias"`      // 微信号
	Remark       string `json:"remark"`     // 备注
	NickName     string `json:"nick_name"`  // 微信名
	LocalType    int    `json:"local_type"` // 2 群聊; 3 群聊成员(非好友); 5,6 企业微信;   1保存的用户和群聊  5的添加的企业微信人员 6群里面的微信人员
	BigHeadUrl   string `json:"big_head_url"`
	SmallHeadUrl string `json:"small_head_url"`
	ExtraBuffer  []byte `json:"extra_buffer"`
}

func (c *ContactModel) Wrap() *model.Contact {
	// buffer := c.ExtraBuffer
	// c.ParseExtraBuffer(buffer)
	return &model.Contact{
		UserName: c.UserName,
		Alias:    c.Alias,
		Remark:   c.Remark,
		NickName: c.NickName,
		IsFriend: c.LocalType != 3,
	}
}

func (c *ContactModel) WrapVo() *model.ContactVo {
	/**
	UserName      string   `json:"userName"`
	Alias         string   `json:"alias"`
	Remark        string   `json:"remark"`
	NickName      string   `json:"nickName"`
	BigHeadUrl    string   `json:"big_head_url"`
	SmallHeadUrl  string   `json:"small_head_url"`
	Labels        []string `json:"labels"`
	Phone         []string `json:"phones"`
	Gender        string   `json:"gender"`
	Signature     string   `json:"signature"`
	Country       string   `json:"country"`
	Province      string   `json:"province"`
	City          string   `json:"city"`
	ChatroomCount int      `json:"chatroom_count"`
	*/
	buffer := c.ExtraBuffer
	cv := &model.ContactVo{
		UserName:     c.UserName,
		Alias:        c.Alias,
		Remark:       c.Remark,
		NickName:     c.NickName,
		BigHeadUrl:   c.BigHeadUrl,
		SmallHeadUrl: c.SmallHeadUrl,
	}
	c.ParseExtraBuffer(buffer, cv)

	return cv
}

func (c *ContactModel) ParseExtraBuffer(b []byte, cv *model.ContactVo) {
	var pbMsg wxproto.ContactInfo
	if err := proto.Unmarshal(b, &pbMsg); err != nil {
		return
	}
	cv.Country = pbMsg.Country
	cv.City = pbMsg.City
	cv.ChatroomCount = int(pbMsg.GetChatroomCount())
	cv.Signature = pbMsg.Signature
	cv.Province = pbMsg.Province

	gender := pbMsg.GetGender()
	switch gender {
	case 1:
		cv.Gender = "男"
	case 2:
		cv.Gender = "女"
	case 0:
		cv.Gender = "未知"
	}
	if pbMsg.PhoneInfo != nil && pbMsg.PhoneInfo.Phones != nil {
		phones := pbMsg.PhoneInfo.Phones
		for _, phone := range phones {
			cv.Phone = append(cv.Phone, phone.PhoneNumer)
		}
	}

	return
}
