package c_home

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"net/http"
)

//后台用户修改自己的密码
func ForgetPassword(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if r.IsMethod("get") {
		return w.
			Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
			View(http.StatusOK, "backend/home/forget.html")
	}

	type Forget struct {
		OldPassword string
		NewPassword string
	}

	mod := &Forget{}
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}

	v := validator.New()
	v.Field("OldPassword").Rule("required", "请输入旧密码")
	v.Field("NewPassword").Rule("password:min=8,max=16",
		"请输入新密码",
		"新密码长度必须在8~16位之间",
		"密码格式有误，请输入数字、字母、符号",
		"密码格式有误，数字、字母、符号至少两种")

	if r, err := v.Validate(mod); err != nil {
		return w.Error(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT, err))
	} else if !r.IsEmpty() {
		return w.Error(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT, r.ToSimpleString()))
	}

	user := s_user.IsSignIn(r.Session())
	if user == nil {
		return w.Error(code.INVALID_ARGUMENT, "错误的登录信息")
	}

	if !s_user.ComparePassword(mod.OldPassword, user.Password) {
		return w.Error(code.INVALID_ARGUMENT, "错误的旧密码")
	}

	var err error
	user.Password, err = s_user.GeneratePassword(mod.NewPassword)
	if err != nil {
		return errors.MarkServer(err)
	}

	if _, err := dao.Dao.ID(user.ID).Update(user); err != nil {
		return errors.MarkServer(err)
	}

	return w.Success()
}
