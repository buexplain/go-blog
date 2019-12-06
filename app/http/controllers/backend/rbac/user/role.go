package c_official_user

import (
	"github.com/buexplain/go-blog/app/boot"
	c_util "github.com/buexplain/go-blog/app/http/controllers/util"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/user"
	"github.com/buexplain/go-blog/services/userRoleRelation"
	"github.com/buexplain/go-fool"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
)

func EditRole(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if !r.IsAjax() {
		user := new(m_user.User)

		user.ID = r.ParamInt("id", 0)
		if user.ID <= 0 {
			return w.JumpBack("参数错误")
		}

		if ok, err := dao.Dao.Get(user); err != nil {
			return ctx.Error().WrapServer(err).Location()
		} else if !ok {
			return w.JumpBack("参数错误")
		}

		if role, err := s_userRoleRelation.GetUserRole(user.ID); err != nil {
			return w.JumpBack(err)
		}else {
			w.Assign("role", template.JS(role.String()))
		}

		return w.
			Assign("user", user).
			Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
			Layout("backend/layout/layout.html").
			View(http.StatusOK, "backend/rbac/user/role.html")
	}

	//开始插入用户角色关系表
	userID := r.ParamInt("id")
	if userID <= 0 {
		return c_util.Error(w, "参数错误")
	}

	roleID := r.FormSliceInt("ids")

	err := s_userRoleRelation.SetUserRole(userID, roleID)
	if err != nil {
		return w.Assign("code", 1).Assign("message", err.Error()).Assign("data", "").JSON(http.StatusOK)
	}

	return c_util.Success(w)
}
