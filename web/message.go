package web

type MessageCode struct {
	Code    int
	Message string
}

var (
	ERROR_SYSTEM        = &MessageCode{500, "系统错误，请稍后重试..."}
	ERROR_FORBID        = &MessageCode{403, "访问无权限"}
	ERROR_MESSAGE       = &MessageCode{999, "%s"}
	ERROR_INVALID_PARAM = &MessageCode{9999, "参数错误，请检查[%s]"}
)
