package c_skeleton

import (
	"github.com/buexplain/go-blog/models/node"
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"html/template"
	"net/http"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	user := s_user.IsSignIn(r.Session())
	menu, err := m_node.GetMenuByUserID(user.ID)
	if err != nil {
		return ctx.Error().WrapServer(err).Location()
	}
	return w.
		Assign("menu", template.JS(menu.String())).
		Assign("user", s_user.IsSignIn(r.Session())).
		View(http.StatusOK, "backend/skeleton/index.html")
}
