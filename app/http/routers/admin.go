package routers

import (
	"github.com/buexplain/go-blog/app/http/controllers/admin/sign"
	"github.com/buexplain/go-blog/app/http/controllers/admin/skeleton"
	"github.com/buexplain/go-blog/app/http/middleware"
	"github.com/buexplain/go-fool"
)

func admin(mux *fool.Mux) {
	// --------------------------登录 开始---------------------------
	mux.Get("/admin/sign", c_sign.Index)
	mux.Post("/admin/sign", c_sign.In)
	mux.Delete("/admin/sign", c_sign.Out)
	// --------------------------登录 结束---------------------------

	mux.Group("", func() {
		mux.Get("/admin/skeleton", c_skeleton.Index)
	}).Use(middleware.IsSignIn)
}
