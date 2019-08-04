package c_skeleton

import (
	"github.com/buexplain/go-blog/app/http/boot/code"
	m_menu "github.com/buexplain/go-blog/models/menu"
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"net/http"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.Assign("user", s_user.IsSignIn(r.Session())).View(http.StatusOK, "backend/skeleton/index.html")
}

//返回所有的菜单
func GetALL(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := m_menu.GetALL()
	if err != nil {
		return ctx.Error().WrapServer(err).Location()
	}
	return w.Assign("data", result).Assign("message", code.Text(code.SUCCESS)).Assign("code", code.SUCCESS).JSON(http.StatusOK)
}
