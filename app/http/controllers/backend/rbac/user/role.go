package c_official_user

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/user"
	"github.com/buexplain/go-blog/services/userRoleRelation"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
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
			return errors.MarkServer(err)
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
		return w.Client(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT))
	}

	roleID := r.FormSliceInt("ids")

	err := s_userRoleRelation.SetUserRole(userID, roleID)
	if err != nil {
		return w.Client(1, err.Error())
	}

	return w.Success()
}
