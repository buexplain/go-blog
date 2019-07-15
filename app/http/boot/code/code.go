package code

//状态码说明
var text map[int]string = map[int]string{}

const StatusTextUnknown = "未知的状态码"

func Text(code int) string {
	if v, ok := text[code]; ok {
		return v
	} else {
		return StatusTextUnknown
	}
}

//正确返回
const SUCCESS = 0

//1 到 99 业务逻辑自定义使用

//100 到 199 预定义客户端错误使用
const CLIENT = 100

//200 到 299 预定义服务端错误使用
const SERVER = 200

//300 以上 服务端预定义异常使用
const EXCEPTION = 300

func init() {
	text[SUCCESS] = "成功"
	text[CLIENT] = "客户端错误"
	text[SERVER] = "服务端错误"
	text[EXCEPTION] = "服务端异常"
}

//预定义的客户端错误
const (
	CODE_INVALID_ROUTE = 101
	CODE_INVALID_AUTH  = 102
	CODE_INVALID_CSRF  = 103
)

func init() {
	text[CODE_INVALID_ROUTE] = "请求地址错误"
	text[CODE_INVALID_AUTH] = "权限校验失败"
	text[CODE_INVALID_CSRF] = "csrf校验失败"
}
