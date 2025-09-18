package types

type WeChatAccountKey struct {
	Account string `json:"account"`
	DataKey string `json:"dataKey"`
	ImgKey  string `json:"imgKey"`
}

type WeChatAccount struct {
	WeChatAccountKey *WeChatAccountKey
	Platform         string
	Version          int
	FullVersion      string
	DataDir          string
	ExePath          string
	Status           string
}
