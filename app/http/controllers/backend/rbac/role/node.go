package c_role

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/role"
	"github.com/buexplain/go-blog/services/roleNodeRelation"
	"github.com/buexplain/go-slim"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
)

func EditNode(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	if !r.IsAjax() {
		role := new(m_role.Role)

		role.ID = r.ParamInt("id", 0)
		if role.ID <= 0 {
			return w.JumpBack(code.Text(code.INVALID_ARGUMENT, "id"))
		}

		if has, err := dao.Dao.Get(role); err != nil {
			return err
		} else if !has {
			return w.JumpBack(code.Text(code.NOT_FOUND_DATA, role.ID))
		}

		if node, err := s_roleNodeRelation.GetRelation(role.ID); err != nil {
			return w.JumpBack(err)
		} else {
			w.Assign("node", template.JS(node.String()))
		}

		return w.
			Assign("role", role).
			Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
			View(http.StatusOK, "backend/rbac/role/node.html")
	}

	//开始插入角色节点关系表
	roleID := r.ParamInt("id")
	if roleID <= 0 {
		return code.NewM(code.INVALID_ARGUMENT, "id")
	}

	nodeID := r.FormSliceInt("ids")

	err := s_roleNodeRelation.SetRelation(roleID, nodeID)
	if err != nil {
		return code.NewM(code.INVALID_ARGUMENT, err)
	}

	return w.Success()
}
