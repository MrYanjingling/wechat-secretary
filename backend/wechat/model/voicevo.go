package model

type VoiceVo struct {
	Key  string `json:"key"`  // MD5
	Data []byte `json:"data"` // for voice
}
