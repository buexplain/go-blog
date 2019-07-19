package c_captcha

import (
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/services/captcha"
	"github.com/buexplain/go-fool"
	"net/http"
)

//显示验证码
func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	url := s_captcha.Generate(r.Session(), s_captcha.SetHeight(20))
	return ctx.Response().
		Assign("code", code.SUCCESS).
		Assign("message", code.Text(code.SUCCESS)).Assign("data", url).
		JSON(http.StatusOK)
}
