package c_role

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	m_role "github.com/buexplain/go-blog/models/role"
	s_role "github.com/buexplain/go-blog/services/role"
	"github.com/buexplain/go-fool"
	"github.com/gorilla/csrf"
	"github.com/thedevsaddam/govalidator"
	"html/template"
	"net/http"
)

//列表
func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := m_role.GetALL()
	if err != nil {
		return ctx.Error().WrapServer(err).Location()
	}
	return w.
		Assign("result", template.JS(result.String())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/rbac/role/index.html")
}

//创建
func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	pid := r.ParamInt("pid", 0)
	return ctx.Response().
		Assign("pid", pid).
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/rbac/role/create.html")
}

//保存
func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	rules := govalidator.MapData{
		"Name": []string{"required"},
	}

	messages := govalidator.MapData{
		"Name": []string{"required:请输入角色名"},
	}

	opts := govalidator.Options{
		Request:         r.Raw(),
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: true,
	}
	v := govalidator.New(opts)
	e := v.Validate()

	if len(e) > 0 {
		return w.JumpBack(e)
	}

	mod := &m_role.Role{}
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}

	if _, err := dao.Dao.Insert(mod); err != nil {
		return w.JumpBack(err)
	}

	return w.JumpBack("操作成功")
}

//编辑
func Edit(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_role.Role)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack("参数错误")
	}

	if has, err := dao.Dao.Get(result); err != nil {
		return ctx.Error().WrapServer(err)
	} else if !has {
		return w.JumpBack("参数错误")
	}

	return w.
		Assign("result", result).
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/rbac/role/create.html")
}

//更新
func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	rules := govalidator.MapData{
		"Name": []string{"required"},
	}

	messages := govalidator.MapData{
		"Name": []string{"required:请输入角色名"},
	}

	opts := govalidator.Options{
		Request:         r.Raw(),
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: true,
	}
	v := govalidator.New(opts)
	e := v.Validate()

	if len(e) > 0 {
		return w.JumpBack(e)
	}

	mod := &m_role.Role{}
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}
	mod.ID = r.ParamInt("id", 0)

	if _, err := dao.Dao.ID(mod.ID).MustCols("Pid").Update(mod); err != nil {
		return w.JumpBack(err)
	}

	return w.Jump("/backend/rbac/role", "操作成功")
}

//删除
func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := r.QuerySliceInt("ids[]")
	if len(ids) == 0 {
		return w.JumpBack("参数错误")
	}
	if _, err := s_role.Destroy(ids); err != nil {
		return w.JumpBack(err)
	}
	return w.Jump("/backend/rbac/role", "操作成功")
}
