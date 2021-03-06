package routers

import (
	"github.com/buexplain/go-blog/app/http/controllers/backend/article/attachment"
	c_category "github.com/buexplain/go-blog/app/http/controllers/backend/article/category"
	c_content "github.com/buexplain/go-blog/app/http/controllers/backend/article/content"
	c_tag "github.com/buexplain/go-blog/app/http/controllers/backend/article/tag"
	c_backup "github.com/buexplain/go-blog/app/http/controllers/backend/backup"
	"github.com/buexplain/go-blog/app/http/controllers/backend/config/group"
	c_configItem "github.com/buexplain/go-blog/app/http/controllers/backend/config/item"
	c_database "github.com/buexplain/go-blog/app/http/controllers/backend/database"
	c_home "github.com/buexplain/go-blog/app/http/controllers/backend/home"
	c_node "github.com/buexplain/go-blog/app/http/controllers/backend/rbac/node"
	c_role "github.com/buexplain/go-blog/app/http/controllers/backend/rbac/role"
	c_official_user "github.com/buexplain/go-blog/app/http/controllers/backend/rbac/user"
	c_profile "github.com/buexplain/go-blog/app/http/controllers/backend/server/profile"
	c_sysInfo "github.com/buexplain/go-blog/app/http/controllers/backend/server/sysInfo"
	c_sysLog "github.com/buexplain/go-blog/app/http/controllers/backend/server/sysLog"
	c_sysReset "github.com/buexplain/go-blog/app/http/controllers/backend/server/sysReset"
	"github.com/buexplain/go-blog/app/http/controllers/backend/sign"
	"github.com/buexplain/go-blog/app/http/controllers/backend/skeleton"
	c_citizen_user "github.com/buexplain/go-blog/app/http/controllers/backend/user"
	"github.com/buexplain/go-blog/app/http/middleware"
	"github.com/buexplain/go-slim"
	"net/http"
	"strings"
)

//后台管理模块路由
func backend(mux *slim.Mux) {
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
		mux.Group("home", func() {
			mux.Get("/", c_home.Index)
			mux.Any("user/forget", c_home.ForgetPassword, http.MethodGet, http.MethodPost)
		})

		//权限管理
		mux.Group("rbac", func() {
			//节点管理
			mux.Restful("node",
				c_node.Index,
				c_node.Create,
				c_node.Store,
				c_node.Edit,
				c_node.Update,
				c_node.Destroy,
			)

			//管理员用户管理
			mux.Restful("user",
				c_official_user.Index,
				c_official_user.Create,
				c_official_user.Store,
				c_official_user.Edit,
				c_official_user.Update,
			)
			mux.Get("user/oauth/:id", c_official_user.OauthList)
			mux.Delete("user/oauth/:id", c_official_user.OauthDestroy)
			mux.Any("user/role/:id", c_official_user.EditRole, http.MethodGet, http.MethodPost)

			//角色管理
			mux.Restful("role",
				c_role.Index,
				c_role.Create,
				c_role.Store,
				c_role.Edit,
				c_role.Update,
				c_role.Destroy,
			)
			mux.Any("role/node/:id", c_role.EditNode, http.MethodGet, http.MethodPost)
		})

		//普通用户管理
		mux.Restful("user",
			c_citizen_user.Index,
			c_citizen_user.Create,
			c_citizen_user.Store,
			c_citizen_user.Edit,
			c_citizen_user.Update,
		)

		//文章管理
		mux.Group("article", func() {
			//分类管理
			mux.Restful("category",
				c_category.Index,
				c_category.Create,
				c_category.Store,
				c_category.Edit,
				c_category.Update,
				c_category.Destroy,
			)

			//标签管理
			mux.Get("tag", c_tag.Index)
			mux.Restful("tag",
				c_tag.Store,
				c_tag.Update,
				c_tag.Destroy,
			)
			mux.Put("tag/delete-batch", c_tag.DestroyBatch)

			//内容管理
			mux.Restful("content",
				c_content.Index,
				c_content.Create,
				c_content.Store,
				c_content.Edit,
				c_content.Update,
				c_content.Destroy,
				c_content.Show,
			)
			mux.Put("content/delete-batch", c_content.DestroyBatch)
			mux.Put("content/online/:id", c_content.Online)
			//返回文章分类
			mux.Get("content/category/:pid", c_content.Category).AddLabel("json")
			//返回文章标签
			mux.Get("content/tag", c_content.Tag).AddLabel("json")
			//添加新的标签
			mux.Post("content/addTag", c_content.AddTag).AddLabel("json")
			//上传文章附件
			mux.Post("content/upload", c_content.Upload).AddLabel("json")

			//附件管理
			mux.Get("attachment", c_attachment.Index)
			mux.Get("attachment/check/:md5", c_attachment.CheckMD5).AddLabel("json")
			mux.Get("attachment/edit/:id", c_attachment.Edit)
			mux.Post("attachment/upload", c_attachment.Upload)
			mux.Get("attachment/download/:id", c_attachment.Download)
			mux.Put("attachment/update/:id", c_attachment.Update).AddLabel("json")
			mux.Delete("attachment/delete/:id", c_attachment.Destroy)
			mux.Put("attachment/delete-batch", c_attachment.DestroyBatch)
		}).Use(middleware.CacheClear)

		//配置管理
		mux.Group("config", func() {
			mux.Restful("group",
				c_configGroup.Index,
				c_configGroup.Create,
				c_configGroup.Store,
				c_configGroup.Edit,
				c_configGroup.Update,
				c_configGroup.Destroy,
			)
			mux.Put("group/delete-batch", c_configGroup.DestroyBatch).AddLabel("json")
			mux.Restful("item",
				c_configItem.Index,
				c_configItem.Create,
				c_configItem.Store,
				c_configItem.Edit,
				c_configItem.Update,
				c_configItem.Destroy,
			)
			mux.Put("item/delete-batch", c_configItem.DestroyBatch).AddLabel("json")
		}).Use(middleware.CacheClear)

		//备份管理
		mux.Get("backup", c_backup.Index)
		mux.Get("backup/start", c_backup.Start).AddLabel("json")
		mux.Get("backup/download", c_backup.Download)
		mux.Delete("backup/delete", c_backup.Destroy)

		//数据管理
		mux.Get("database", c_database.Index)
		mux.Post("database", c_database.SQL)

		//服务器管理
		mux.Group("server", func() {
			//服务器状态
			mux.Get("info", c_sysInfo.Index)
			//进程信息
			mux.Get("debug/pprof/:name", c_profile.Index).Use(middleware.RbacCheck).Use(func(ctx *slim.Ctx, w *slim.Response, r *slim.Request) {
				r.Raw().URL.Path = strings.TrimPrefix(r.Raw().URL.Path, "/backend/server")
				ctx.Next()
			})
			//服务器重启
			mux.Group("reset", func() {
				mux.Get("index", c_sysReset.Index)
				mux.Post("start", c_sysReset.Start)
				mux.Get("check", c_sysReset.Check)
			})
			//系统日志
			mux.Group("sysLog", func() {
				mux.Get("", c_sysLog.Index)
				mux.Get("download", c_sysLog.Download)
				mux.Get("show", c_sysLog.Show)
				mux.Delete("delete", c_sysLog.Destroy)
			})
		})
	}).Use(middleware.RbacCheck)
	// --------------------------需要权限校验的路由 结束---------------------------
}
