package c_configGroup

import (
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	m_configGroup "github.com/buexplain/go-blog/models/config/group"
	s_configGroup "github.com/buexplain/go-blog/services/config/group"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"net/http"
)


//表单校验器
var v *validator.Validator

func init() {
	v = validator.New()
	v.Field("Name").Rule("required", "请输入组名称")
	v.Field("key").Rule("required", "请输入组字段名称")
}

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if !r.IsAjax() {
		return w.Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
			View(http.StatusOK, "backend/config/group/index.html")
	}
	counter, result, err := s_configGroup.GetList(ctx)
	if err != nil {
		return err
	}
	w.Assign("count", counter)
	return w.Success(result)
}

func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.
		Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		View(http.StatusOK, "backend/config/group/create.html")
}


func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := new(m_configGroup.ConfigGroup)
	if err := r.FormToStruct(mod); err != nil {
		return err
	}

	if r, err := v.Validate(mod); err != nil {
		return errors.MarkServer(err)
	} else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	if _, err := dao.Dao.Insert(mod); err != nil {
		return errors.MarkServer(err)
	}

	return w.JumpBack("操作成功")
}

func Edit(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_configGroup.ConfigGroup)

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
		View(http.StatusOK, "backend/config/group/create.html")
}

func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := new(m_configGroup.ConfigGroup)
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}
	mod.ID = r.ParamInt("id", 0)

	vClone := v.Clone()
	vClone.Field("ID").Rule("required", "ID错误")

	if r, err := vClone.Validate(mod); err != nil {
		return errors.MarkServer(err)
	} else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	if _, err := dao.Dao.ID(mod.ID).Update(mod); err != nil {
		return errors.MarkServer(err)
	}

	return w.Jump("/backend/config/group", "操作成功")
}

func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := []int{r.ParamInt("id", 0)}
	_, err := s_configGroup.Destroy(ids)
	if err != nil {
		return err
	}
	return w.Success()
}

//批量删除
func DestroyBatch(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := r.FormSliceInt("ids")
	_, err := s_configGroup.Destroy(ids)
	if err != nil {
		return err
	}
	return w.Success()
}