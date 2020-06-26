package c_home

import (
	"github.com/buexplain/go-blog/services/user"
	s_userRoleRelation "github.com/buexplain/go-blog/services/userRoleRelation"
	"github.com/buexplain/go-slim"
	"net/http"
)

func Index(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	user := s_user.IsSignIn(r)
	userRoleName, err := s_userRoleRelation.GetRoleNameByUserID(user.ID)
	if err != nil {
		return err
	}
	w.Assign("userRoleName", userRoleName)
	return w.Assign("user", user).
		View(http.StatusOK, "backend/home/index.html")
}
