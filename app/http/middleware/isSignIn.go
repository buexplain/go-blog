package middleware

import (
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"net/http"
)

//判断用户是否登录了
func IsSignIn(ctx *fool.Ctx, w *fool.Response, r *fool.Request) {
	if u := s_user.IsSignIn(r.Session()); u != nil {
		ctx.Next()
	} else {
		if ctx.Route() != nil && ctx.Route().HasLabel("json") {
			//存在路由，并且路由有json标签，则响应json格式
			ctx.Throw(ctx.Response().Assign("code", code.CODE_INVALID_AUTH).Assign("message", code.Text(code.CODE_INVALID_AUTH)).Assign("data", "").JSON(http.StatusOK))
		} else {
			ctx.Throw(ctx.Response().Abort(http.StatusFound, code.Text(code.CODE_INVALID_AUTH)))
		}
	}
}
