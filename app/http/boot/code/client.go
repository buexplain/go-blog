package code

import "github.com/buexplain/go-slim/errors"

const (
	//1 ~ 400 客户端错误
	CLIENT        = errors.ClientCode
	INVALID_ROUTE = iota + CLIENT - 100
	INVALID_AUTH
	INVALID_CSRF
	INVALID_ARGUMENT
	NOT_FOUND_DATA
)

func init() {
	text[CLIENT] = "客户端错误"
	text[INVALID_ROUTE] = "请求地址错误"
	text[INVALID_AUTH] = "权限校验失败"
	text[INVALID_CSRF] = "csrf校验失败"
	text[INVALID_ARGUMENT] = "错误的参数"
	text[NOT_FOUND_DATA] = "没有找到相关数据"
}
