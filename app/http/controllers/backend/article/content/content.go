package c_content

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	m_category "github.com/buexplain/go-blog/models/category"
	m_content "github.com/buexplain/go-blog/models/content"
	m_contentTag "github.com/buexplain/go-blog/models/contentTag"
	"github.com/buexplain/go-blog/models/tag"
	"github.com/buexplain/go-blog/models/util"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"net/http"
	"path/filepath"
	"strconv"
)

//表单校验器
var v *validator.Validator

func init()  {
	v = validator.New()
	v.Rule("Title").Add("required", "请填写标题")
	v.Rule("Category").Add("required", "请选择分类")
	//校验标签是否存在
	v.Custom("CheckUnique", func(field string, value interface{}, rule *validator.Rule, structVar interface{}) (s string, e error) {
		str, ok := value.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		if !m_util.CheckUnique("Tag", field, str, rule.GetInt("id")) {
			return rule.Message(0), nil
		}
		return "", nil
	})
}

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	query := m_util.NewQuery("Content", ctx).Limit()

	query.Finder.Desc("ID")

	var result m_content.List
	query.Find(&result)

	count := query.Count()

	if query.Error != nil {
		return ctx.Error().WrapServer(query.Error).Location()
	}

	return w.Assign("count", count).
		Assign("result", result).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/article/content/index.html")
}

func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	tagList := new(m_tag.List)
	if err := dao.Dao.Find(tagList); err != nil {
		return err
	}
	w.Assign("tagList", tagList)

	return w.
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").View(http.StatusOK, "backend/article/content/create.html")
}

func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := new(m_tag.Tag)
	if err := r.FormToStruct(mod); err != nil {
		return err
	}

	if r, err := v.Validate(mod); err != nil {
		return ctx.Error().WrapServer(err)
	}else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	if _, err := dao.Dao.Insert(mod); err != nil {
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

	tagList := new(m_tag.List)
	if err := dao.Dao.Find(tagList); err != nil {
		return ctx.Error().WrapServer(err)
	}
	w.Assign("tagList", tagList)

	contentTagList := make(m_contentTag.List, 0)
	if err := dao.Dao.Where("ContentID=?", result.ID).Find(&contentTagList); err == nil {
		w.Assign("contentTagList", contentTagList)
	} else {
		return ctx.Error().WrapServer(err)
	}

	return w.
		Assign("result", result).
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/article/content/create.html")
}

func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := new(m_tag.Tag)
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}
	mod.ID = r.ParamInt("id", 0)

	vClone := v.Clone()
	vClone.Rule("ID").Add("required", "ID错误")
	vClone.Rule("Name").Add("CheckUnique:id="+strconv.Itoa(mod.ID), "该标签名已存在")

	if r, err := vClone.Validate(mod); err != nil {
		return ctx.Error().WrapServer(err)
	}else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	if _, err := dao.Dao.ID(mod.ID).Update(mod); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.Jump("/backend/article/content", "操作成功")
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

	return w.Jump("/backend/article/content", "操作成功")
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

//上传图片
func Upload(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	uploads, err := r.FileSlice("file[]")
	if err != nil {
		return ctx.Error().WrapClient(err)
	}
	defer func() {
		_ = uploads.Close()
	}()
	uploads.SetValidateExt(boot.Config.Business.Upload.Ext...)
	_, err = uploads.SaveToPath(filepath.Join(boot.ROOT_PATH, boot.Config.Business.Upload.Save))
	if err != nil {
		return ctx.Error().WrapClient(err)
	}
	data := []map[string]string{}
	for _, upload := range uploads  {
		data = append(data, map[string]string{"name":upload.Name(), "path":"/"+upload.Result()})
	}
	return w.Assign("data", data).
		Assign("message", code.Text(code.SUCCESS)).
		Assign("code", code.SUCCESS).
		JSON(http.StatusOK)
}
