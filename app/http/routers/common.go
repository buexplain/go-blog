package routers

import (
	"github.com/buexplain/go-blog/app/http/controllers/common/captcha"
	"github.com/buexplain/go-fool"
)

func common(mux *fool.Mux)  {
	mux.Get("/common/captcha", c_captcha.Index)
}