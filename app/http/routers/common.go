package routers

import (
	"github.com/buexplain/go-blog/app/http/controllers/common/captcha"
	"github.com/buexplain/go-blog/app/http/controllers/common/jump"
	c_oAuth "github.com/buexplain/go-blog/app/http/controllers/common/oAuth"
	"github.com/buexplain/go-slim"
)

func common(mux *slim.Mux) {
	//图形验证码输出
	mux.Get("/common/captcha", c_captcha.Index)
	mux.Any("/common/jump", c_jump.Index)
	mux.Get("/common/oauth", c_oAuth.Index)
	mux.Post("/common/oauth/register", c_oAuth.Register).AddLabel("json")
	mux.Post("/common/oauth/bind", c_oAuth.Bind).AddLabel("json")
}
