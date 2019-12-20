package c_content

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/category"
	"github.com/buexplain/go-blog/models/content"
	"github.com/buexplain/go-blog/models/tag"
	"github.com/buexplain/go-blog/services"
	"github.com/buexplain/go-blog/services/attachment"
	"github.com/buexplain/go-blog/services/content"
	"github.com/buexplain/go-blog/services/tag"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"net/http"
	"xorm.io/xorm"
)

//表单校验器
var v *validator.Validator

//初始化表单校验器
func init() {
	v = validator.New()
	v.Field("Title").Rule("required", "请填写标题")
	v.Field("Category").Rule("required", "请选择分类")
	v.Field("Online").Rule(fmt.Sprintf("in:in=%d,%d", m_content.OnlineYes, m_content.OnlineNo), "请选择上下线")
	v.Field("Body").Rule("required", "请填写内容")
}

//列表
func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if !r.IsAjax() {
		return w.
			Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
			Layout("backend/layout/layout.html").
			View(http.StatusOK, "backend/article/content/index.html")
	}
	query := s_services.NewQuery("Content", ctx).Limit().Where()
	query.Finder.Desc("ID")
	var result m_content.List
	var count int64
	query.FindAndCount(&result, &count)
	if query.Error != nil {
		return ctx.Error().WrapServer(query.Error).Location()
	}
	return w.
		Assign("code", code.SUCCESS).
		Assign("message", "操作成功").
		Assign("count", count).
		Assign("data", result).
		JSON(http.StatusOK)
}

//新增
func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	tagList := new(m_tag.List)
	if err := dao.Dao.Find(tagList); err != nil {
		return err
	}
	w.Assign("tagList", tagList)

	return w.
		Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").View(http.StatusOK, "backend/article/content/create.html")
}

//保存
func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := new(m_content.Content)
	if err := r.FormToStruct(mod); err != nil {
		return err
	}

	if r, err := v.Validate(mod); err != nil {
		return ctx.Error().WrapServer(err)
	}else if !r.IsEmpty() {
		return w.Assign("data", "").
			Assign("message", r.ToSimpleString()).
			Assign("code", 1).
			JSON(http.StatusOK)
	}

	tagsID := r.FormSliceInt("tagsID")

	if len(tagsID) == 0 {
		return w.Assign("data", "").
			Assign("message", "请选择标签").
			Assign("code", 1).
			JSON(http.StatusOK)
	}

	if err := s_content.Save(mod, tagsID, 0); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.Assign("data", mod.ID).
		Assign("message", code.Text(code.SUCCESS)).
		Assign("code", code.SUCCESS).
		JSON(http.StatusOK)
}

//编辑
func Edit(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	//内容
	result := new(m_content.Content)
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
		Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").View(http.StatusOK, "backend/article/content/create.html")
}

//更新
func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := new(m_content.Content)
	if err := r.FormToStruct(mod); err != nil {
		return err
	}
	mod.ID = r.ParamInt("id", 0)
	vClone := v.Clone()
	vClone.Field("ID").Rule("required", "ID错误")

	if r, err := vClone.Validate(mod); err != nil {
		return ctx.Error().WrapServer(err)
	}else if !r.IsEmpty() {
		return w.Assign("data", "").
			Assign("message", r.ToSimpleString()).
			Assign("code", 1).
			JSON(http.StatusOK)
	}

	tagsID := r.FormSliceInt("tagsID")

	if len(tagsID) == 0 {
		return w.Assign("data", "").
			Assign("message", "请选择标签").
			Assign("code", 1).
			JSON(http.StatusOK)
	}

	if err := s_content.Save(mod, tagsID, mod.ID); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.Assign("data", mod.ID).
		Assign("message", code.Text(code.SUCCESS)).
		Assign("code", code.SUCCESS).
		JSON(http.StatusOK)
}

//单个删除
func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_content.Content)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack("参数错误")
	}

	if affected, err := dao.Dao.Where("Online=?", m_content.OnlineNo).Delete(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}else if affected > 0 {
		return w.
			Assign("code", code.SUCCESS).
			Assign("message", "操作成功").
			Assign("data", "").
			JSON(http.StatusOK)
	}

	return w.
		Assign("code", 1).
		Assign("message", "操作失败").
		Assign("data", "").
		JSON(http.StatusOK)
}

//批量删除
func DestroyBatch(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := r.FormSliceInt("ids")
	if len(ids) == 0 {
		return w.JumpBack("参数错误")
	}

	if affected, err := s_services.DestroyBatch("Content", ids, func(session *xorm.Session) *xorm.Session {
		return  session.Where("Online=?", m_content.OnlineNo)
	}); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}else if affected > 0 {
		return w.
			Assign("code", code.SUCCESS).
			Assign("message", "操作成功").
			Assign("data", "").
			JSON(http.StatusOK)
	}

	return w.
		Assign("code", 1).
		Assign("message", "操作失败").
		Assign("data", "").
		JSON(http.StatusOK)
}

//查看
func Show(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	d, err := s_content.GetDetails(r.ParamInt("id"))
	if err != nil {
		return w.Assign("data", d).Assign("code", 1).Assign("message", err.Error()).JSON(http.StatusOK)
	}
	return w.Assign("data", d).Assign("code", code.SUCCESS).Assign("message", code.Text(code.SUCCESS)).JSON(http.StatusOK)
}

//设置上下线
func Online(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_content.Content)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack("参数错误")
	}
	result.Online = r.FormInt("online", 0)

	if result.Online == m_content.OnlineYes {
		result.Online = m_content.OnlineNo
	}else {
		result.Online = m_content.OnlineYes
	}

	if _, err := dao.Dao.ID(result.ID).Update(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.
		Assign("code", code.SUCCESS).
		Assign("message", "操作成功").
		Assign("data", result.Online).
		JSON(http.StatusOK)
}

//返回分类
func Category(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	pid := r.ParamInt("pid", -1)
	query := dao.Dao.Table("Category").Desc("ID")
	if pid > -1 {
		query.Where("Pid=?", pid)
	}
	result := make(m_category.List, 0)
	if err := query.Find(&result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}
	return w.Assign("data", result).
		Assign("message", code.Text(code.SUCCESS)).
		Assign("code", code.SUCCESS).
		JSON(http.StatusOK)
}

//返回标签
func Tag(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_tag.List)
	if err := dao.Dao.Find(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}
	return w.Assign("data", result).
		Assign("message", code.Text(code.SUCCESS)).
		Assign("code", code.SUCCESS).
		JSON(http.StatusOK)
}

//新增tag
func AddTag(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	name := r.Form("name", "")
	id, err := s_tag.Store(name)
	if err != nil {
		return w.Assign("data", "").
			Assign("message", err.Error()).
			Assign("code", code.SUCCESS).
			JSON(http.StatusOK)
	}
	return w.Assign("data", id).
		Assign("message", code.Text(code.SUCCESS)).
		Assign("code", code.SUCCESS).
		JSON(http.StatusOK)
}

//上传附件
func Upload(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	file, err := r.File("file")
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	result, err := s_attachment.Upload(file, "")
	if err != nil {
		return w.Assign("data", "").Assign("message", err.Error()).Assign("code", code.SERVER).JSON(http.StatusOK)
	}
	return w.Assign("data", result).Assign("message", code.Text(code.SUCCESS)).Assign("code", code.SUCCESS).JSON(http.StatusOK)
}
