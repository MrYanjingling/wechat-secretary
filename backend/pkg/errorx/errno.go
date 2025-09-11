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
	ErrWriteConflict           = 100000020
	ErrInternal                = 100000021
	ErrWeChatNameNotFound      = 100000022
	ErrWeChatProcessNotExist   = 100000023
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

func WriteConflict() error {
	return New(ErrWriteConflict, KV("msg", "Write conflict"))
}

func Internal() error {
	return New(ErrInternal, KV("msg", "internal error"))
}

func WeChatAccountNotFound(name string) error {
	return New(ErrWeChatNameNotFound, KVf("msg", "WeChat %s name account not found", name))
}

func WeChatProcessNotExist() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func PlatformUnsupported(platform string, version int) error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func ErrWeChatOffline() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func ErrSIPEnabled() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func ErrValidatorNotSet() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func ErrNoValidKey() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func ErrWeChatDLLNotFound() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func OpenProcessFailed() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func RunCmdFailed() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func ErrDecryptHashVerificationFailed() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func DecryptCreateCipherFailed() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func ErrAlreadyDecrypted() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func IncompleteRead() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func StatFileFailed(path string) error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func OpenFileFailed(path string) error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func ReadFileFailed(path string) error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func ErrNoMemoryRegionsFound() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func CreatePipeFileFailed() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

// func RunCmdFailed() error {
// 	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
// }

func ReadPipeFileFailed() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func OpenPipeFileFailed() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func ReadMemoryFailed() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func ErrReadMemoryTimeout() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func DecodeKeyFailed() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func WriteOutputFailed() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func ErrDecryptIncorrectKey() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}

func ErrDecryptOperationCanceled() error {
	return New(ErrWeChatProcessNotExist, KVf("msg", "WeChat process not exist"))
}
