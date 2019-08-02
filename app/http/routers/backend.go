package routers

import (
	"github.com/buexplain/go-blog/app/http/controllers/backend/home"
	"github.com/buexplain/go-blog/app/http/controllers/backend/menu"
	"github.com/buexplain/go-blog/app/http/controllers/backend/sign"
	"github.com/buexplain/go-blog/app/http/controllers/backend/skeleton"
	"github.com/buexplain/go-blog/app/http/middleware"
	"github.com/buexplain/go-fool"
)

func backend(mux *fool.Mux) {
	// --------------------------登录 开始---------------------------
	mux.Get("/backend/sign", c_sign.Index)
	mux.Post("/backend/sign", c_sign.In)
	mux.Delete("/backend/sign", c_sign.Out)
	// --------------------------登录 结束---------------------------

	// --------------------------需要权限校验的路由 开始---------------------------
	mux.Group("", func() {
		mux.Get("/backend/skeleton", c_skeleton.Index)
		mux.Get("/backend/home", c_home.Index)
		mux.Get("/backend/menu", c_menu.Index)
		mux.Get("/backend/menu/list", c_menu.GetList).AddLabel("json")
		mux.Get("/backend/menu/create/:pid", c_menu.Create)
		mux.Post("/backend/menu", c_menu.Store)
		mux.Delete("/backend/menu", c_menu.Destroy)

	}).Use(middleware.IsSignIn)
	// --------------------------需要权限校验的路由 结束---------------------------
}
