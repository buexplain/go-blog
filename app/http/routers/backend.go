package routers

import (
	c_category "github.com/buexplain/go-blog/app/http/controllers/backend/article/category"
	c_tag "github.com/buexplain/go-blog/app/http/controllers/backend/article/tag"
	"github.com/buexplain/go-blog/app/http/controllers/backend/home"
	c_node "github.com/buexplain/go-blog/app/http/controllers/backend/rbac/node"
	c_role "github.com/buexplain/go-blog/app/http/controllers/backend/rbac/role"
	c_user "github.com/buexplain/go-blog/app/http/controllers/backend/rbac/user"
	"github.com/buexplain/go-blog/app/http/controllers/backend/sign"
	"github.com/buexplain/go-blog/app/http/controllers/backend/skeleton"
	"github.com/buexplain/go-blog/app/http/middleware"
	"github.com/buexplain/go-fool"
	"net/http"
)

//后台管理模块路由
func backend(mux *fool.Mux) {
	// --------------------------登录 开始---------------------------
	mux.Get("/backend/sign", c_sign.Index)
	mux.Post("/backend/sign", c_sign.In)
	mux.Delete("/backend/sign", c_sign.Out)
	// --------------------------登录 结束---------------------------

	// --------------------------需要权限校验的路由 开始---------------------------
	mux.Group("/backend", func() {
		//后台骨架
		mux.Get("skeleton", c_skeleton.Index)

		//我的桌面
		mux.Get("home", c_home.Index)

		//权限管理
		mux.Group("rbac", func() {
			//节点管理
			mux.Get("node", c_node.Index)
			mux.Get("node/create/:pid", c_node.Create)
			mux.Post("node", c_node.Store)
			mux.Get("node/edit/:id", c_node.Edit)
			mux.Put("node/:id", c_node.Update)
			mux.Delete("node", c_node.Destroy)

			//用户管理
			mux.Any("user/forget", c_user.Forget, http.MethodGet, http.MethodPost)
			mux.Get("user", c_user.Index)
			mux.Get("user/create", c_user.Create)
			mux.Post("user", c_user.Store)
			mux.Get("user/edit/:id", c_user.Edit)
			mux.Put("user/:id", c_user.Update)
			mux.Any("user/role/:id", c_user.EditRole,  http.MethodGet, http.MethodPost)

			//角色管理
			mux.Get("role", c_role.Index)
			mux.Get("role/create/:pid", c_role.Create)
			mux.Post("role", c_role.Store)
			mux.Get("role/edit/:id", c_role.Edit)
			mux.Put("role/:id", c_role.Update)
			mux.Delete("role", c_role.Destroy)
		})

		//文章管理
		mux.Group("article", func() {
			//分类管理
			mux.Get("category", c_category.Index)
			mux.Get("category/create/:pid", c_category.Create)
			mux.Post("category", c_category.Store)
			mux.Get("category/edit/:id", c_category.Edit)
			mux.Put("category/:id", c_category.Update)
			mux.Delete("category", c_category.Destroy)

			//标签管理
			mux.Get("tag", c_tag.Index)
			mux.Get("tag/create", c_tag.Create)
			mux.Post("tag", c_tag.Store)
			mux.Get("tag/edit/:id", c_tag.Edit)
			mux.Put("tag/:id", c_tag.Update)
			mux.Delete("tag/:id", c_tag.Destroy)
		})

	}).Use(middleware.IsSignIn)
	// --------------------------需要权限校验的路由 结束---------------------------
}
