package c_sign

import (
	"fmt"
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/services/captcha"
	s_oauth "github.com/buexplain/go-blog/services/oauth"
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-slim"
	"github.com/buexplain/go-slim/errors"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"net/http"
	"strings"
)

//表单校验器
var v *validator.Validator

func init() {
	v = validator.New()
	v.Field("Account").Rule("required", "请输入账号")
	v.Field("Password").Rule("required", "请输入密码")
	v.Field("CaptchaVal").Rule("VerifyCaptcha", "请输入验证码", "验证码错误")
}

//显示登录页面
func Index(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	if s_user.IsSignIn(r) != nil {
		return w.Redirect(http.StatusFound, "/backend/skeleton")
	}
	w.Assign("github", s_oauth.NewGithub().GetURL("user", "/backend/skeleton", r))
	return w.Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).View(http.StatusOK, "backend/sign/index.html")
}

//登录
func In(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	type In struct {
		Account    string
		Password   string
		CaptchaVal string
	}
	mod := &In{}
	if err := r.FormToStruct(mod); err != nil {
		return errors.MarkClient(err)
	}

	vClone := v.Clone()
	vClone.Custom("VerifyCaptcha", func(field string, value interface{}, rule *validator.Rule, structVar interface{}) (s string, e error) {
		str, ok := value.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		str = strings.TrimSpace(str)
		if str == "" {
			return rule.Message(0), nil
		}
		if s_captcha.Verify(r.Session(), str) == false {
			return rule.Message(1), nil
		}
		return "", nil
	})

	if r, err := vClone.Validate(mod); err != nil {
		return err
	} else if !r.IsEmpty() {
		return w.Error(code.INVALID_ARGUMENT, r.ToSimpleString())
	}

	var err error
	_, err = s_user.OfficialSignIn(ctx.Request().Session(), mod.Account, mod.Password)
	if err != nil {
		return err
	}
	return w.Success()
}

//退出登录
func Out(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	s_user.SignOut(r.Session())
	redirect := r.Query("redirect", "/backend/sign")
	return w.Redirect(http.StatusFound, redirect)
}
