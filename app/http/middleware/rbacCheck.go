package middleware

import (
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/models/user"
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
	}else {
		if rbacCheck(ctx) {
			ctx.Next()
		} else {
			if ctx.Route() != nil && ctx.Route().HasLabel("json") {
				//存在路由，并且路由有json标签，则响应json格式
				ctx.Throw(ctx.Response().Assign("code", code.INVALID_AUTH).Assign("message", code.Text(code.INVALID_AUTH)).Assign("data", "").JSON(http.StatusOK))
			} else {
				ctx.Throw(ctx.Response().Jump("/backend/sign", code.Text(code.INVALID_AUTH)))
			}
		}
	}
}

func rbacCheck(ctx *fool.Ctx) bool {
	return true
	user := s_user.IsSignIn(ctx.Request().Session())

	//判断后台用户是否登录
	if user == nil || user.Identity != m_user.IdentityOfficial {
		return false
	}

	//得到当前命中的路由
	route := ctx.Route()

	//得到当前路由所在角色的列表
	nodeRoleIDList, err := s_roleNodeRelation.GetRoleIDByNodeURL(route.GetPath())
	if err != nil {
		ctx.Logger().ErrorF("得到当前路由所在角色的列表错误: %s", ctx.Error().WrapServer(err).Location())
		return false
	}

	//得到当前用户拥有的角色
	userRoleIDList, err := s_userRoleRelation.GetRoleIDByUserID(user.ID)
	if err != nil {
		ctx.Logger().ErrorF("得到当前用户拥有的角色错误: %s", ctx.Error().WrapServer(err).Location())
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