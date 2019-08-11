package c_tag

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	m_tag "github.com/buexplain/go-blog/models/tag"
	m_util "github.com/buexplain/go-blog/models/util"
	"github.com/buexplain/go-fool"
	"github.com/gorilla/csrf"
	"net/http"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	query := m_util.NewQuery("Tag", ctx).Limit()

	query.Finder.Desc("ID")

	var result m_tag.List
	query.Find(&result)

	count := query.Count()

	if query.Error != nil {
		return ctx.Error().WrapServer(query.Error).Location()
	}

	return w.Assign("count", count).
		Assign("result", result).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/article/tag/index.html")
}

func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").View(http.StatusOK, "backend/article/tag/create.html")
}

func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_tag.Tag)

	if err := r.FormToStruct(result); err != nil {
		return err
	}

	if result.Name == "" {
		return w.JumpBack("请填写标签名")
	}

	if !m_util.CheckUnique("Tag", "Name", result.Name) {
		return w.JumpBack("该标签名已存在")
	}

	if _, err := dao.Dao.Insert(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.JumpBack("操作成功")
}

func Edit(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_tag.Tag)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack("参数错误")
	}

	if ok, err := dao.Dao.Get(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	} else if !ok {
		return w.JumpBack("参数错误")
	}

	return w.
		Assign("result", result).
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/article/tag/create.html")
}

func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_tag.Tag)

	if err := r.FormToStruct(result); err != nil {
		return w.JumpBack(err)
	}

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack("ID错误")
	}

	if result.Name == "" {
		return w.JumpBack("请填写标签名")
	}

	if !m_util.CheckUnique("Tag", "Name", result.Name, result.ID) {
		return w.JumpBack("该标签名已存在")
	}

	if _, err := dao.Dao.ID(result.ID).Update(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.Jump("/backend/article/tag", "操作成功")
}

func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_tag.Tag)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack("参数错误")
	}

	if _, err := dao.Dao.Delete(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.Jump("/backend/article/tag", "操作成功")
}
