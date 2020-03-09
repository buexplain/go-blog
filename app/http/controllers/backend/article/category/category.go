package c_category

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/category"
	"github.com/buexplain/go-blog/services/category"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
)

//表单校验器
var v *validator.Validator

func init() {
	v = validator.New()
	v.Field("Name").Rule("required", "请输入分类名")
}

//列表
func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := s_category.GetALL()
	return w.
		Assign("result", template.JS(result.String())).
		View(http.StatusOK, "backend/article/category/index.html")
}

//创建
func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	pid := r.QueryInt("pid", 0)
	return ctx.Response().
		Assign("pid", pid).
		Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		View(http.StatusOK, "backend/article/category/create.html")
}

//保存
func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := &m_category.Category{}
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}

	if r, err := v.Validate(mod); err != nil {
		return errors.MarkClient(err)
	} else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	if _, err := s_category.Store(mod); err != nil {
		return w.JumpBack(err)
	}

	return w.JumpBack("操作成功")
}

//编辑
func Edit(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_category.Category)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack(code.Text(code.INVALID_ARGUMENT, "id"))
	}

	if has, err := dao.Dao.Get(result); err != nil {
		return errors.MarkServer(err)
	} else if !has {
		return w.JumpBack(code.Text(code.NOT_FOUND_DATA, result.ID))
	}

	return w.
		Assign("result", result).
		Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		View(http.StatusOK, "backend/article/category/create.html")
}

//更新
func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := &m_category.Category{}
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}
	mod.ID = r.ParamInt("id", 0)

	vClone := v.Clone()
	vClone.Field("ID").Rule("required", "ID错误")

	if r, err := v.Validate(mod); err != nil {
		return errors.MarkServer(err)
	} else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	if _, err := s_category.Store(mod); err != nil {
		return w.JumpBack(err)
	}

	return w.Jump("/backend/article/category", "操作成功")
}

//删除
func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := r.QuerySlicePositiveInt("ids")
	if len(ids) == 0 {
		return w.Jump("/backend/article/category", "入参错误")
	}
	if affected, err := s_category.Destroy(ids); err != nil {
		return w.Jump("/backend/article/category", err)
	} else {
		return w.Jump("/backend/article/category", fmt.Sprintf("操作 %d 条数据成功", affected))
	}
}
