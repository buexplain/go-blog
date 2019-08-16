package c_role

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/role"
	"github.com/buexplain/go-blog/services/roleNodeRelation"
	"github.com/buexplain/go-fool"
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
			return ctx.Error().WrapServer(err).Location()
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
			Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
			Layout("backend/layout/layout.html").
			View(http.StatusOK, "backend/rbac/role/node.html")
	}

	//开始插入角色节点关系表
	roleID := r.ParamInt("id")
	if roleID <= 0 {
		return w.Assign("code", 1).Assign("message", "参数错误").Assign("data", "").JSON(http.StatusOK)
	}

	nodeID := r.FormSliceInt("ids[]")

	err := s_roleNodeRelation.SetRoleNode(roleID, nodeID)
	if err != nil {
		return w.Assign("code", 1).Assign("message", err.Error()).Assign("data", "").JSON(http.StatusOK)
	}

	return w.Assign("code", 0).Assign("message", "操作成功").Assign("data", "").JSON(http.StatusOK)
}

