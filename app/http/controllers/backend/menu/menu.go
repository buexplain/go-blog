package c_menu

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/menu"
	"github.com/buexplain/go-fool"
	"github.com/gorilla/csrf"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/menu/index.html")
}

func GetList(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := m_menu.GetList()
	if err != nil {
		return ctx.Error().WrapServer(err).Location()
	}
	return w.Assign("data", result).Assign("message", code.Text(code.SUCCESS)).Assign("code", code.SUCCESS).JSON(http.StatusOK)
}

func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	pid := r.ParamInt("pid", 0)
	return ctx.Response().
		Assign("pid", pid).
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/menu/create.html")
}

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
		return ctx.Error().WrapClient(err).Location()
	}

	if _, err := dao.Dao.Insert(mod); err != nil {
		return err
	}

	return w.Redirect(http.StatusFound, "/backend/menu")
}

func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := r.QuerySliceInt("ids[]")
	if len(ids) == 0 {
		return w.Assign("message", "参数错误").Assign("code", 1).JSON(http.StatusOK)
	}
	if _, err := dao.Dao.In("ID", ids).Delete(new(m_menu.Menu)); err != nil {
		return ctx.Error().WrapServer(err)
	}
	return w.Assign("data", "").Assign("message", code.Text(code.SUCCESS)).Assign("code", code.SUCCESS).JSON(http.StatusOK)
}
