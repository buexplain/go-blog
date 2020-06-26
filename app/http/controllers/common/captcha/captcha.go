package c_captcha

import (
	"github.com/buexplain/go-blog/services/captcha"
	"github.com/buexplain/go-slim"
)

//显示验证码
func Index(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	url := s_captcha.Generate(r.Session(), 38, 104, 4)
	return w.Success(url)
}
