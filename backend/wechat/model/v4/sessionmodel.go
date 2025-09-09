package v4

import (
	"time"
	"wechat-secretary/backend/wechat/model"
)

type SessionModel struct {
	Username              string `json:"username"`
	Summary               string `json:"summary"`
	LastTimestamp         int    `json:"last_timestamp"`
	LastMsgSender         string `json:"last_msg_sender"`
	LastSenderDisplayName string `json:"last_sender_display_name"`

	// Type                     int    `json:"type"`
	// UnreadCount              int    `json:"unread_count"`
	// UnreadFirstMsgSrvID      int    `json:"unread_first_msg_srv_id"`
	// IsHidden                 int    `json:"is_hidden"`
	// Draft                    string `json:"draft"`
	// Status                   int    `json:"status"`
	// SortTimestamp            int    `json:"sort_timestamp"`
	// LastClearUnreadTimestamp int    `json:"last_clear_unread_timestamp"`
	// LastMsgLocaldID          int    `json:"last_msg_locald_id"`
	// LastMsgType              int    `json:"last_msg_type"`
	// LastMsgSubType           int    `json:"last_msg_sub_type"`
	// LastMsgExtType           int    `json:"last_msg_ext_type"`
}

func (s *SessionModel) Wrap() *model.SessionVo {
	return &model.SessionVo{
		UserName: s.Username,
		Order:    s.LastTimestamp,
		// NickName: s.LastSenderDisplayName,
		Content: s.Summary,
		Time:    time.Unix(int64(s.LastTimestamp), 0),
	}
}
