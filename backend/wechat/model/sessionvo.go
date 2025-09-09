package model

import (
	"time"
)

type SessionVo struct {
	UserName string    `json:"userName"`
	Order    int       `json:"order"`
	NickName string    `json:"nickName"`
	Content  string    `json:"content"`
	HeadUrl  string    `json:"headUrl"`
	Time     time.Time `json:"time"`
}

// CREATE TABLE Session(
// strUsrName TEXT  PRIMARY KEY,
// nOrder INT DEFAULT 0,
// nUnReadCount INTEGER DEFAULT 0,
// parentRef TEXT,
// Reserved0 INTEGER DEFAULT 0,
// Reserved1 TEXT,
// strNickName TEXT,
// nStatus INTEGER,
// nIsSend INTEGER,
// strContent TEXT,
// nMsgType	INTEGER,
// nMsgLocalID INTEGER,
// nMsgStatus INTEGER,
// nTime INTEGER,
// editContent TEXT,
// othersAtMe INT,
// Reserved2 INTEGER DEFAULT 0,
// Reserved3 TEXT,
// Reserved4 INTEGER DEFAULT 0,
// Reserved5 TEXT,
// bytesXml BLOB
// )
