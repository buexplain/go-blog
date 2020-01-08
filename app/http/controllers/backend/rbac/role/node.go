package c_role

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/role"
	"github.com/buexplain/go-blog/services/roleNodeRelation"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
)

func EditNode(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if !r.IsAjax() {
		role := new(m_role.Role)

		role.ID = r.ParamInt("id", 0)
		if role.ID <= 0 {
			return w.JumpBack("参数错误")
		}

		if ok, err := dao.Dao.Get(role); err != nil {
			return errors.MarkServer(err)
		} else if !ok {
			return w.JumpBack("参数错误")
		}

		if node, err := s_roleNodeRelation.GetRoleNode(role.ID); err != nil {
			return w.JumpBack(err)
		}else {
			w.Assign("node", template.JS(node.String()))
		}

		return w.
			Assign("role", role).
			Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
			Layout("backend/layout/layout.html").
			View(http.StatusOK, "backend/rbac/role/node.html")
	}

	//开始插入角色节点关系表
	roleID := r.ParamInt("id")
	if roleID <= 0 {
		return w.Client(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT))
	}

	nodeID := r.FormSliceInt("ids")

	err := s_roleNodeRelation.SetRoleNode(roleID, nodeID)
	if err != nil {
		return w.Client(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT), err.Error())
	}

	return w.Success()
}

