package types

type WeChatAccountKey struct {
	Account string `json:"account"`
	DataKey string `json:"dataKey"`
	ImgKey  string `json:"imgKey"`
	DataDir string `json:"dataDir"`
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

type WeChatAccountDetails struct {
	Account        string `json:"account"`
	DataKey        string `json:"dataKey"`
	ImgKey         string `json:"imgKey"`
	DataDecryptDir string `json:"dataDir"`
	Platform       string `json:"platform"`
	Version        int    `json:"version"`
	FullVersion    string `json:"fullVersion"`
	DataSourceDir  string `json:"dataSourceDir"`
	ExePath        string `json:"exePath"`
	NonLocal       bool   `json:"nonLocal"`
	Status         string
}
