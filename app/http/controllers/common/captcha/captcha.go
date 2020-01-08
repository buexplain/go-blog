package c_captcha

import (
	"github.com/buexplain/go-blog/services/captcha"
	"github.com/buexplain/go-fool"
)

//显示验证码
func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	url := s_captcha.Generate(r.Session())
	return w.Success(url)
}
