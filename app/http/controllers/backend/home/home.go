package c_home

import (
	"github.com/buexplain/go-blog/services/user"
	s_userRoleRelation "github.com/buexplain/go-blog/services/userRoleRelation"
	"github.com/buexplain/go-fool"
	"net/http"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	user := s_user.IsSignIn(r)
	userRoleName, err := s_userRoleRelation.GetRoleNameByUserID(user.ID)
	if err != nil {
		return err
	}
	w.Assign("userRoleName", userRoleName)
	return w.Assign("user", user).
		View(http.StatusOK, "backend/home/index.html")
}
