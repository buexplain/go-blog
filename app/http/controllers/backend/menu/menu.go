package c_menu

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/menu"
	s_menu "github.com/buexplain/go-blog/services/menu"
	"github.com/buexplain/go-fool"
	"github.com/gorilla/csrf"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

//列表
func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/menu/index.html")
}

//返回所有的菜单
func GetALL(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := m_menu.GetALL()
	if err != nil {
		return ctx.Error().WrapServer(err).Location()
	}
	return w.Assign("data", result).Assign("message", code.Text(code.SUCCESS)).Assign("code", code.SUCCESS).JSON(http.StatusOK)
}

//创建
func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	pid := r.ParamInt("pid", 0)
	return ctx.Response().
		Assign("pid", pid).
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/menu/create.html")
}

//保存
func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	rules := govalidator.MapData{
		"Name": []string{"required"},
		"URL": []string{"required"},
	}

	messages := govalidator.MapData{
		"Name": []string{"required:请输入菜单名"},
		"URL": []string{"required:请输入访问路径"},
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
		return w.Assign("message", e).Jump(r.Raw().URL.Path, e)
	}

	mod := &m_menu.Menu{}
	if err := r.FormToStruct(mod); err != nil {
		return w.Jump(r.Raw().Header.Get("Referer"), ctx.Error().WrapClient(err))
	}

	if _, err := dao.Dao.Insert(mod); err != nil {
		return w.Jump("/backend/menu", err)
	}

	return w.Redirect(http.StatusFound, "/backend/menu")
}

//编辑
func Edit(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_menu.Menu)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.Abort(http.StatusBadRequest, "入参错误")
	}

	if has, err := dao.Dao.Get(result); err != nil {
		return ctx.Error().WrapServer(err)
	} else if !has {
		return w.Abort(http.StatusBadRequest, "ID错误")
	}

	return ctx.Response().
		Assign("result", result).
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/menu/create.html")
}

//更新
func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	rules := govalidator.MapData{
		"Name": []string{"required"},
		"URL": []string{"required"},
	}

	messages := govalidator.MapData{
		"Name": []string{"required:请输入菜单名"},
		"URL": []string{"required:请输入访问路径"},
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
		return w.Assign("message", e).Jump(r.Raw().URL.Path, e)
	}

	mod := &m_menu.Menu{}
	if err := r.FormToStruct(mod); err != nil {
		return w.Jump(r.Raw().Header.Get("Referer"), ctx.Error().WrapClient(err))
	}
	mod.ID = r.ParamInt("id", 0)

	if _, err := dao.Dao.ID(mod.ID).MustCols("Pid").Update(mod); err != nil {
		return w.Jump("/backend/menu", err)
	}

	return w.Redirect(http.StatusFound, "/backend/menu")
}

//删除
func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := r.QuerySliceInt("ids[]")
	if len(ids) == 0 {
		return w.Jump("/backend/menu", "入参错误")
	}
	if _, err := s_menu.Destroy(ids); err != nil {
		return w.Jump("/backend/menu", err)
	}
	return w.Redirect(http.StatusFound, "/backend/menu")
}
