package middleware

import (
	h_boot "github.com/buexplain/go-blog/app/http/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/models/user"
	s_node "github.com/buexplain/go-blog/services/node"
	"github.com/buexplain/go-blog/services/roleNodeRelation"
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-blog/services/userRoleRelation"
	"github.com/buexplain/go-slim"
)

//rbac 权限校验
func RbacCheck(ctx *slim.Ctx, w *slim.Response, r *slim.Request) {
	if ctx.Route() == nil {
		//该中间件必须设置为组中间件或者是路由中间件
		ctx.Throw(code.New(code.MIDDLEWARE_SET_ERROR))
		return
	}
	if rbacCheck(ctx) {
		ctx.Next()
		return
	}
	node := s_node.GetByURL(ctx.Route().GetPath(), ctx.Request().Raw().Method)
	if node == nil {
		ctx.Throw(code.NewF(code.SERVER, "路由 %s %s 没有写入到Node表", ctx.Request().Raw().Method, ctx.Route().GetPath()))
	} else {
		ctx.Throw(code.New(code.INVALID_AUTH))
	}
}

func rbacCheck(ctx *slim.Ctx) bool {
	user := s_user.IsSignIn(ctx.Request())

	//判断后台用户是否登录
	if user == nil || user.Identity != m_user.IdentityOfficial {
		return false
	}

	//得到当前命中的路由
	route := ctx.Route()

	//得到当前路由所在角色的列表
	nodeRoleIDList, err := s_roleNodeRelation.GetRoleIDByNodeURL(route.GetPath(), ctx.Request().Raw().Method)
	if err != nil {
		h_boot.Logger.ErrorF("获取当前路由所在角色的列表错误: %s", err)
		return false
	}

	//得到当前用户拥有的角色
	var userRoleIDList s_userRoleRelation.UserRoleIDList
	userRoleIDList, err = s_userRoleRelation.GetRoleIDByUserID(user.ID)
	if err != nil {
		h_boot.Logger.ErrorF("获取当前用户拥有的角色错误: %s", err)
		return false
	}

	//计算交集
	for _, userRoleID := range userRoleIDList {
		for _, nodeRoleID := range nodeRoleIDList {
			//存在交集，权限校验通过
			if int(userRoleID) == int(nodeRoleID) {
				return true
			}
		}
	}

	//权限校验失败
	return false
}
