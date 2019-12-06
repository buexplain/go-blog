package code

//1 到 99 业务逻辑自定义使用
//100 到 499 预定义客户端错误使用
const CLIENT = 100
func init() {
	text[CLIENT] = "客户端错误"
}

//预定义的客户端错误
const (
	INVALID_ROUTE = 101
	INVALID_AUTH  = 102
	INVALID_CSRF  = 103
	INVALID_ARGUMENT = 104
)

func init() {
	text[INVALID_ROUTE] = "请求地址错误"
	text[INVALID_AUTH] = "权限校验失败"
	text[INVALID_CSRF] = "csrf校验失败"
	text[INVALID_ARGUMENT] = "错误的参数"
}
