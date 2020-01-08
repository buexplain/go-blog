package c_sign

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/services/captcha"
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"net/http"
	"strings"
)


//显示登录页面
func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if s_user.IsSignIn(r.Session()) != nil {
		return w.Redirect(http.StatusFound, "/backend/skeleton")
	}
	return w.Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).View(http.StatusOK, "backend/sign/index.html")
}

//登录
func In(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	type In struct {
		Account string
		Password string
		CaptchaVal string
	}
	mod := &In{}
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}

	v := validator.New()
	v.Field("Account").Rule("required", "请输入账号")
	v.Field("Password").Rule("required", "请输入密码")
	v.Field("CaptchaVal").Rule("VerifyCaptcha", "请输入验证码", "验证码错误")
	v.Custom("VerifyCaptcha", func(field string, value interface{}, rule *validator.Rule, structVar interface{}) (s string, e error) {
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

	if r, err := v.Validate(mod); err != nil {
		return errors.MarkServer(err)
	}else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	var err error
	_, err = s_user.OfficialSignIn(ctx.Request().Session(), mod.Account, mod.Password)
	if err != nil {
		return w.JumpBack(err)
	}

	return w.Redirect(http.StatusFound, "/backend/skeleton")
}

//退出登录
func Out(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	s_user.SignOut(r.Session())
	return w.Redirect(http.StatusFound, "/backend/sign")
}
