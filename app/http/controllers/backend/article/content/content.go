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
	"github.com/buexplain/go-fool/errors"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"net/http"
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
		return w.Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
			View(http.StatusOK, "backend/article/content/index.html")
	}
	query := s_services.NewQuery("Content", ctx).Limit().Where()
	query.Finder.Desc("ID")
	var result m_content.List
	var count int64
	query.FindAndCount(&result, &count)
	if query.Error != nil {
		return errors.MarkServer(query.Error)
	}
	w.Assign("count", count)
	return w.Success(result)
}

//新增
func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	tagList := new(m_tag.List)
	if err := dao.Dao.Find(tagList); err != nil {
		return err
	}
	w.Assign("tagList", tagList)

	return w.Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		View(http.StatusOK, "backend/article/content/create.html")
}

//保存
func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := new(m_content.Content)
	if err := r.FormToStruct(mod); err != nil {
		return err
	}

	if r, err := v.Validate(mod); err != nil {
		return err
	} else if !r.IsEmpty() {
		return w.Error(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT, r.ToSimpleString()))
	}

	tagsID := r.FormSliceInt("tagsID")

	if len(tagsID) == 0 {
		return w.Error(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT, "id"))
	}

	if err := s_content.Save(mod, tagsID, 0); err != nil {
		return err
	}

	return w.Success(mod.ID)
}

//编辑
func Edit(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	//内容
	result := new(m_content.Content)
	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack(code.Text(code.INVALID_ARGUMENT, "id"))
	}
	if has, err := dao.Dao.Get(result); err != nil {
		return errors.MarkServer(err)
	} else if !has {
		return w.JumpBack(code.Text(code.NOT_FOUND_DATA, result.ID))
	}

	return w.Assign("result", result).
		Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		View(http.StatusOK, "backend/article/content/create.html")
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
		return err
	} else if !r.IsEmpty() {
		return w.Error(code.INVALID_ARGUMENT, r.ToSimpleString())
	}

	tagsID := r.FormSliceInt("tagsID")

	if len(tagsID) == 0 {
		return w.Error(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT, "id"))
	}

	if err := s_content.Save(mod, tagsID, mod.ID); err != nil {
		return err
	}

	return w.Success(mod.ID)
}

//单个删除
func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := []int{r.ParamInt("id", 0)}
	err := s_content.DestroyBatch(ids)
	if err != nil {
		return err
	}
	return w.Success()
}

//批量删除
func DestroyBatch(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := r.FormSliceInt("ids")
	err := s_content.DestroyBatch(ids)
	if err != nil {
		return err
	}
	return w.Success()
}

//查看
func Show(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := s_content.GetDetails(r.ParamInt("id"))
	if r.IsAjax() {
		if err != nil {
			return err
		}
		return w.Success(result)
	}
	if err != nil {
		return err
	}
	return w.Assign("result", result).
		View(http.StatusOK, "backend/article/content/show.html")
}

//设置上下线
func Online(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_content.Content)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack(code.Text(code.INVALID_ARGUMENT, "id"))
	}
	result.Online = r.FormInt("online", 0)

	if result.Online == m_content.OnlineYes {
		result.Online = m_content.OnlineNo
	} else {
		result.Online = m_content.OnlineYes
	}

	if _, err := dao.Dao.ID(result.ID).Update(result); err != nil {
		return errors.MarkServer(err)
	}

	return w.Success(result.Online)
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
		return err
	}
	return w.Success(result)
}

//返回标签
func Tag(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_tag.List)
	if err := dao.Dao.Find(result); err != nil {
		return errors.MarkServer(err)
	}
	return w.Success(result)
}

//新增tag
func AddTag(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	name := r.Form("name", "")
	id, err := s_tag.Store(name)
	if err != nil {
		return err
	}
	return w.Success(id)
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
		return err
	}
	return w.Success(result)
}

func Render(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	var err error
	var html string
	id := r.QueryInt("id", 0)
	if id > 0 {
		html, err = s_content.RenderByID(id)
	} else {
		markdown := r.Form("markdown", "")
		if markdown != "" {
			html, err = s_content.Render(markdown)
		}
	}
	if err != nil {
		return w.Error(code.CLIENT, err.Error())
	}
	return w.Success(html)
}
