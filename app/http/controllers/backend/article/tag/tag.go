package c_tag

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/tag"
	"github.com/buexplain/go-blog/services"
	s_tag "github.com/buexplain/go-blog/services/tag"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"net/http"
	"strconv"
)

//表单校验器
var v *validator.Validator

func init() {
	v = validator.New()
	v.Field("Name").Rule("required", "请填写标签名").Rule("CheckUnique:id=0", "该标签名已存在")
	//校验标签是否存在
	v.Custom("CheckUnique", func(field string, value interface{}, rule *validator.Rule, structVar interface{}) (s string, e error) {
		str, ok := value.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		if !s_services.CheckUnique("Tag", field, str, rule.GetInt("id")) {
			return rule.Message(0), nil
		}
		return "", nil
	})
}

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if !r.IsAjax() {
		return w.Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
			View(http.StatusOK, "backend/article/tag/index.html")
	}
	counter, result, err := s_tag.GetList(ctx)
	if err != nil {
		return err
	}
	w.Assign("count", counter)
	return w.Success(result)
}

func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	names := r.FormSliceByComma("names")
	if len(names) == 0 {
		return w.Error(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT, "names"))
	}
	_, err := s_tag.Stores(names)
	if err != nil {
		return w.JumpBack(err.Error())
	}
	return w.RedirectBack()
}

func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := new(m_tag.Tag)
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}
	mod.ID = r.ParamInt("id", 0)

	vClone := v.Clone()
	vClone.Field("ID").Rule("required", "ID错误")
	vClone.Field("Name").Rule("CheckUnique:id="+strconv.Itoa(mod.ID), "该标签名已存在")

	if r, err := vClone.Validate(mod); err != nil {
		return errors.MarkServer(err)
	} else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	if _, err := dao.Dao.ID(mod.ID).Update(mod); err != nil {
		return err
	}

	return w.Success()
}

//单个删除
func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := []int{r.ParamInt("id", 0)}
	_, err := s_tag.Destroy(ids)
	if err != nil {
		return err
	}
	return w.Success()
}

//批量删除
func DestroyBatch(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := r.FormSliceInt("ids")
	_, err := s_tag.Destroy(ids)
	if err != nil {
		return err
	}
	return w.Success()
}
