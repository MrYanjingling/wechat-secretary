package service

import "time"

type Config struct {
	DataKey  string    `json:"dataKey"`
	DataDir  string    `json:"dataDir"`
	WorkDir  string    `json:"workDir"`
	Platform string    `json:"platform"`
	Version  string    `json:"version"`
	Time     time.Time `json:"time"`
}

type CoreService struct {
	config *Config
}
