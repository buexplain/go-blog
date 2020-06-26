package c_skeleton

import (
	"github.com/buexplain/go-blog/models/node"
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-slim"
	"html/template"
	"net/http"
)

func Index(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	user := s_user.IsSignIn(r)
	menu, err := m_node.GetMenuByUserID(user.ID)
	if err != nil {
		return err
	}
	return w.
		Assign("menu", template.JS(menu.String())).
		Assign("user", s_user.IsSignIn(r)).
		View(http.StatusOK, "backend/skeleton/index.html")
}
