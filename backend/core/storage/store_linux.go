package storage

const (
	storePath = "/var/lib/WechatSecretary"
)

func isEphemeralError(err error) bool {
	return false
}
