package c_home

import (
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"net/http"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.Assign("user", s_user.IsSignIn(r.Session())).
		View(http.StatusOK, "backend/home/index.html")
}
