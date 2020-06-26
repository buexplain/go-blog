package routers

import "github.com/buexplain/go-slim"

//设置全局路由正则
func regexp(mux *slim.Mux) {
	mux.Regexp("id", `^[1-9][0-9]*$`)
}
