package model

import (
	"google.golang.org/protobuf/proto"
	"wechat-secretary/backend/wechatdb/model/wxproto"
)

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
type ContactV4 struct {
	UserName     string `json:"username"`   // 微信号id  gz_ 公众号
	Alias        string `json:"alias"`      // 微信号
	Remark       string `json:"remark"`     // 备注
	NickName     string `json:"nick_name"`  // 微信名
	LocalType    int    `json:"local_type"` // 2 群聊; 3 群聊成员(非好友); 5,6 企业微信;   1保存的用户和群聊  5的添加的企业微信人员 6群里面的微信人员
	BigHeadUrl   string `json:"big_head_url"`
	SmallHeadUrl string `json:"small_head_url"`
	ExtraBuffer  []byte `json:"extra_buffer"`
}

func (c *ContactV4) Wrap() *Contact {
	buffer := c.ExtraBuffer
	c.ParseExtraBuffer(buffer)
	return &Contact{
		UserName: c.UserName,
		Alias:    c.Alias,
		Remark:   c.Remark,
		NickName: c.NickName,
		IsFriend: c.LocalType != 3,
	}
}
func (c *ContactV4) ParseExtraBuffer(b []byte) {
	var pbMsg wxproto.ContactInfo
	if err := proto.Unmarshal(b, &pbMsg); err != nil {
		return
	}

	return
}
