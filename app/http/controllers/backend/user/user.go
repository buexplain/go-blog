package c_citizen_user

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/user"
	"github.com/buexplain/go-blog/services"
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"net/http"
	"strconv"
)


//表单校验器
var v *validator.Validator

func init()  {
	v = validator.New()
	v.Field("Account").Rule("required", "请填写账号").Rule("CheckUnique:id=0", "该账号已存在")
	v.Field("Password",).Rule("password:min=8,max=16",
		"请输入新密码",
		"新密码长度必须在8~16位之间",
		"密码格式有误，请输入数字、字母、符号",
		"密码格式有误，数字、字母、符号至少两种")
	v.Field("Status").Rule("in:in=1,2", "请选择状态")
	//校验账号是否存在
	v.Custom("CheckUnique", func(field string, value interface{}, rule *validator.Rule, structVar interface{}) (s string, e error) {
		str, ok := value.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		if !s_services.CheckUnique("User", field, str, rule.GetInt("id")) {
			return rule.Message(0), nil
		}
		return "", nil
	})
}

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	query := s_services.NewQuery("User", ctx).Limit()
	//此处只展示普通用户
	query.Finder.Where("Identity=?", m_user.IdentityCitizen)
	query.Finder.Desc("ID")
	query.Where()
	var result m_user.List
	var count int64
	query.FindAndCount(&result, &count)
	if query.Error != nil {
		return errors.MarkServer(query.Error)
	}
	return w.
		Assign("count", count).
		Assign("result", result).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/user/index.html")
}

func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.
		Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").View(http.StatusOK, "backend/user/create.html")
}

func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := new(m_user.User)
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}

	if r, err := v.Validate(mod); err != nil {
		return errors.MarkServer(err)
	}else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	if mod.Nickname == "" {
		mod.Nickname = mod.Account
	}

	if p, err := s_user.GeneratePassword(mod.Password); err != nil {
		return w.JumpBack(err)
	}else {
		mod.Password = p
	}

	//强制用户身份为普通用户
	mod.Identity = m_user.IdentityCitizen

	if _, err := dao.Dao.Insert(mod); err != nil {
		return errors.MarkServer(err)
	}

	return w.JumpBack("操作成功")
}

func Edit(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_user.User)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack("参数错误")
	}

	if ok, err := dao.Dao.Where("Identity=?", m_user.IdentityCitizen).Get(result); err != nil {
		return errors.MarkServer(err)
	} else if !ok {
		return w.JumpBack("参数错误")
	}

	return w.
		Assign("result", result).
		Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/user/create.html")
}

func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := new(m_user.User)
	if err := r.FormToStruct(mod); err != nil {
		return err
	}
	mod.ID = r.ParamInt("id", 0)

	vClone := v.Clone()
	vClone.Custom("password", func(field string, value interface{}, rule *validator.Rule, structVar interface{}) (s string, e error) {
		str, _ := value.(string)
		if str == "" {
			return "", nil
		}
		return validator.Pool("password")(field, value, rule, structVar)
	})
	vClone.Field("ID").Rule("required", "ID错误")
	vClone.Field("Account").Rule("CheckUnique:id="+strconv.Itoa(mod.ID), "该账号已存在")

	if r, err := vClone.Validate(mod); err != nil {
		return errors.MarkServer(err)
	}else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	if mod.Nickname == "" {
		mod.Nickname = mod.Account
	}

	if mod.Password != "" {
		if p, err := s_user.GeneratePassword(mod.Password); err != nil {
			return w.JumpBack(err)
		}else {
			mod.Password = p
		}
	}

	//强制用户身份为普通用户
	mod.Identity = m_user.IdentityCitizen

	//强制只允许修改非管理员用户
	if _, err := dao.Dao.ID(mod.ID).Where("Identity=?", m_user.IdentityCitizen).Update(mod); err != nil {
		return errors.MarkServer(err)
	}

	return w.Jump("/backend/user", "操作成功")
}