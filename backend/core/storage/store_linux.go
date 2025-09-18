package storage

const (
	storePath             = "/var/lib/WechatSecretary"
	WechatSecretaryPrefix = "/var/wechatSecretary/data"
)

func isEphemeralError(err error) bool {
	return false
}
