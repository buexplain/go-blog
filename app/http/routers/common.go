package routers

import (
	"github.com/buexplain/go-blog/app/http/controllers/common/captcha"
	"github.com/buexplain/go-blog/app/http/controllers/common/jump"
	"github.com/buexplain/go-fool"
)

func common(mux *fool.Mux) {
	//图形验证码输出
	mux.Get("/common/captcha", c_captcha.Index)
	mux.Any("/common/jump", c_jump.Index)
}
