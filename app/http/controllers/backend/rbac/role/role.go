package c_role

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	m_role "github.com/buexplain/go-blog/models/role"
	s_role "github.com/buexplain/go-blog/services/role"
	"github.com/buexplain/go-slim"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
)

//表单校验器
var v *validator.Validator

func init() {
	v = validator.New()
	v.Field("Name").Rule("required", "请输入角色名")
}

//列表
func Index(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	result, err := m_role.GetALL()
	if err != nil {
		return err
	}
	return w.
		Assign("result", template.JS(result.String())).
		View(http.StatusOK, "backend/rbac/role/index.html")
}

//创建
func Create(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	pid := r.QueryInt("pid", 0)
	return ctx.Response().
		Assign("pid", pid).
		Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		View(http.StatusOK, "backend/rbac/role/create.html")
}

//保存
func Store(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	mod := &m_role.Role{}
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}

	if r, err := v.Validate(mod); err != nil {
		return err
	} else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	if _, err := dao.Dao.Insert(mod); err != nil {
		return w.JumpBack(err)
	}

	return w.JumpBack("操作成功")
}

//编辑
func Edit(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	result := new(m_role.Role)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack(code.Text(code.INVALID_ARGUMENT, "id"))
	}

	if has, err := dao.Dao.Get(result); err != nil {
		return err
	} else if !has {
		return w.JumpBack(code.Text(code.NOT_FOUND_DATA, result.ID))
	}

	return w.
		Assign("result", result).
		Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		View(http.StatusOK, "backend/rbac/role/create.html")
}

//更新
func Update(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	mod := &m_role.Role{}
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}
	mod.ID = r.ParamInt("id", 0)

	vClone := v.Clone()
	vClone.Field("ID").Rule("required", "ID错误")

	if r, err := vClone.Validate(mod); err != nil {
		return err
	} else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	if _, err := dao.Dao.ID(mod.ID).MustCols("Pid").Update(mod); err != nil {
		return w.JumpBack(err)
	}

	return w.Jump("/backend/rbac/role", "操作成功")
}

//删除
func Destroy(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	ids := r.QuerySlicePositiveInt("ids")
	if len(ids) == 0 {
		return w.JumpBack(code.Text(code.INVALID_ARGUMENT, "ids"))
	}
	if affected, err := s_role.Destroy(ids); err != nil {
		return w.JumpBack(err)
	} else {
		return w.Jump("/backend/rbac/role", fmt.Sprintf("操作 %d 条数据成功", affected))
	}
}
