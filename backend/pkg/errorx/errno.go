package errorx

const (
	// ErrAgentInvalidParamCode               = 100000000
	// ErrAgentSupportedChatModelProtocol     = 100000001
	// ErrAgentResourceNotFound               = 100000002
	// ErrAgentPermissionCode                 = 100000003
	// ErrAgentIDGenFailCode                  = 100000004
	// ErrAgentCreateDraftCode                = 100000005
	// ErrAgentGetCode                        = 100000006
	// ErrAgentUpdateCode                     = 100000007
	// ErrAgentSetDraftBotDisplayInfo         = 100000008
	// ErrAgentGetDraftBotDisplayInfoNotFound = 100000009
	// ErrAgentPublishSingleAgentCode         = 100000010
	// ErrAgentAlreadyBindDatabaseCode        = 100000011
	// ErrAgentExecuteErrCode                 = 100000012
	// ErrAgentNoModelInUseCode               = 100000013
	ErrFileGroupNotExist       = 100000014
	ErrWechatDbFileNotExist    = 100000015
	ErrMessageDbFileInitFailed = 100000016
	ErrQueryFailed             = 100000017
	ErrMediaTypeUnsupported    = 100000018
	ErrKeyEmpty                = 100000019
)

func FileGroupNotFound(name string) error {
	return New(ErrFileGroupNotExist, KVf("msg", "File group %s not exists", name))
}

func DBFileNotFound(path, pattern string) error {
	return New(ErrWechatDbFileNotExist, KVf("msg", "Wechat db file not exist %s:%s", path, pattern))
}

func MessageDBInitFailed() error {
	return New(ErrMessageDbFileInitFailed, KVf("msg", "Failed to init message db"))
}

func QueryFailed() error {
	return New(ErrQueryFailed, KV("msg", "Failed to query data"))
}

func MediaTypeUnsupported(_type string) error {
	return New(ErrMediaTypeUnsupported, KVf("msg", "Unsupported media type %s", _type))
}

func KeyEmpty() error {
	return New(ErrQueryFailed, KV("msg", "Empty key"))
}
