package storage

import (
	"errors"
	"github.com/labstack/gommon/log"
	"golang.org/x/sys/windows"
	"os/user"
	"path/filepath"
	"syscall"
)

var (
	storePath = getStorePath()
)

func getStorePath() string {
	if u, err := user.Current(); err == nil {
		return filepath.Join(u.HomeDir, "WechatSecretary")
	} else {
		log.Errorf("Failed to get home dir: %s", err)
		return "./WechatSecretary"
	}
}

func isEphemeralError(err error) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		switch errno {
		case windows.ERROR_SHARING_VIOLATION:
			return true
		}
	}
	return false
}
