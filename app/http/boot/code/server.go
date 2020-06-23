package code

import "github.com/buexplain/go-fool/errors"

//预定义的服务端错误
const (
	// >= 500 服务端错误
	SERVER           = errors.ServerCode
	CALL_THIRD_ERROR = iota + SERVER
	MIDDLEWARE_SET_ERROR
	INVALID_CONFIG
)

func init() {
	text[SERVER] = "服务端错误"
	text[CALL_THIRD_ERROR] = "第三方接口错误"
	text[MIDDLEWARE_SET_ERROR] = "中间件设置错误"
	text[INVALID_CONFIG] = "配置错误"
}
