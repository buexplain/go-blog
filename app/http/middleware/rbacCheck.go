package middleware

import (
	"fmt"
	h_boot "github.com/buexplain/go-blog/app/http/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/models/user"
	s_node "github.com/buexplain/go-blog/services/node"
	"github.com/buexplain/go-blog/services/roleNodeRelation"
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-blog/services/userRoleRelation"
	"github.com/buexplain/go-fool"
	"net/http"
)

//rbac 权限校验
func RbacCheck(ctx *fool.Ctx, w *fool.Response, r *fool.Request) {
	if ctx.Route() == nil {
		//该中间件必须设置为组中间件或者是路由中间件
		ctx.Throw(ctx.Response().Plain(http.StatusInternalServerError, code.Text(code.MIDDLEWARE_SET_ERROR)))
	} else {
		if rbacCheck(ctx) {
			ctx.Next()
		} else {
			var message string
			if ctx.App().Debug() && ctx.Route() != nil {
				node := s_node.GetByURL(ctx.Route().GetPath(), ctx.Request().Raw().Method)
				if node == nil {
					h_boot.Logger.WarningF("路由 %s 没有写入到Node表", ctx.Route().GetPath())
					message = fmt.Sprintf("%s: %s", code.Text(code.INVALID_AUTH), ctx.Route().GetPath())
				}else {
					message = fmt.Sprintf("%s: %s %s %s", code.Text(code.INVALID_AUTH), node.Name, node.URL, node.Methods)
				}
			} else {
				message = code.Text(code.INVALID_AUTH)
			}
			if ctx.Route() != nil && ctx.Route().HasLabel("json") {
				//存在路由，并且路由有json标签，则响应json格式
				ctx.Throw(w.Error(code.INVALID_AUTH, message))
			} else {
				ctx.Throw(ctx.Response().JumpBack(message))
			}
		}
	}
}

func rbacCheck(ctx *fool.Ctx) bool {
	user := s_user.IsSignIn(ctx.Request().Session())

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
	userRoleIDList, err := s_userRoleRelation.GetRoleIDByUserID(user.ID)
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
