package v4

import (
	"path/filepath"
	"wechat-secretary/backend/wechatdb/model"
)

type MediaModel struct {
	Type       string `json:"type"`
	Key        string `json:"key"`
	Dir1       string `json:"dir1"`
	Dir2       string `json:"dir2"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	ModifyTime int64  `json:"modifyTime"`
}

func (m *MediaModel) Wrap() *model.MediaVo {

	var path string
	switch m.Type {
	case "image":
		path = filepath.Join("msg", "attach", m.Dir1, m.Dir2, "Img", m.Name)
	case "video":
		path = filepath.Join("msg", "video", m.Dir1, m.Name)
	case "file":
		path = filepath.Join("msg", "file", m.Dir1, m.Name)
	}

	return &model.MediaVo{
		Type:       m.Type,
		Key:        m.Key,
		Path:       path,
		Name:       m.Name,
		Size:       m.Size,
		ModifyTime: m.ModifyTime,
	}
}
