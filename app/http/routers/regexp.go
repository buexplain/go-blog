package routers

import "github.com/buexplain/go-fool"

//设置全局路由正则
func regexp(mux *fool.Mux) {
	mux.Regexp("id", `^[1-9][0-9]*$`)
}
