package routers

import (
	"github.com/buexplain/go-blog/app/http/controllers/backend/home"
	"github.com/buexplain/go-blog/app/http/controllers/backend/menu"
	"github.com/buexplain/go-blog/app/http/controllers/backend/sign"
	"github.com/buexplain/go-blog/app/http/controllers/backend/skeleton"
	c_user "github.com/buexplain/go-blog/app/http/controllers/backend/user"
	"github.com/buexplain/go-blog/app/http/middleware"
	"github.com/buexplain/go-fool"
	"net/http"
)

func backend(mux *fool.Mux) {
	// --------------------------登录 开始---------------------------
	mux.Get("/backend/sign", c_sign.Index)
	mux.Post("/backend/sign", c_sign.In)
	mux.Delete("/backend/sign", c_sign.Out)
	// --------------------------登录 结束---------------------------

	// --------------------------需要权限校验的路由 开始---------------------------
	mux.Group("/backend", func() {
		mux.Get("skeleton", c_skeleton.Index)
		mux.Get("skeleton/all", c_skeleton.GetALL).AddLabel("json")
		mux.Get("home", c_home.Index)
		mux.Get("menu", c_menu.Index)
		mux.Get("menu/all", c_menu.GetALL).AddLabel("json")
		mux.Get("menu/create/:pid", c_menu.Create)
		mux.Post("menu", c_menu.Store)
		mux.Get("menu/edit/:id", c_menu.Edit)
		mux.Put("menu/:id", c_menu.Update)
		mux.Delete("menu", c_menu.Destroy)
		mux.Any("user/forget", c_user.Forget, http.MethodGet, http.MethodPut)
	}).Use(middleware.IsSignIn)
	// --------------------------需要权限校验的路由 结束---------------------------
}
