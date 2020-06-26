package c_jump

import (
	"github.com/buexplain/go-slim"
)

//页面跳转显示信息
func Index(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	message := r.Input("message")
	url := r.Input("url")
	wait := r.InputInt("wait", 5)
	return w.Jump(url, message, wait)
}
