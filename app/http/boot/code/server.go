package code

//500以上预定义服务端错误使用
const SERVER = 500

func init() {
	text[SERVER] = "服务端错误"
}

//预定义的服务端错误
const (
	MIDDLEWARE_SET_ERROR = 501
)

func init() {
	text[MIDDLEWARE_SET_ERROR] = "中间件设置错误"
}
