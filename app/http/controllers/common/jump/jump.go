package c_jump

import (
	"github.com/buexplain/go-fool"
)

//页面跳转显示信息
func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	message := r.Input("message")
	url := r.Input("url")
	wait := r.InputInt("wait", 5)
	return w.Jump(url, message, wait)
}
