package c_sign

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/services/captcha"
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"github.com/gorilla/csrf"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"strings"
)

//显示登录页面
func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if s_user.IsSignIn(r.Session()) != nil {
		return w.Redirect(http.StatusFound, "/backend/skeleton")
	}
	return w.Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).View(http.StatusOK, "backend/sign/index.html")
}

//登录
func In(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if s_captcha.Verify(r.Session(), strings.TrimSpace(r.Form("captchaVal", ""))) == false {
		return w.Jump(r.Raw().URL.Path, "验证码错误")
	}

	rules := govalidator.MapData{
		"account":  []string{"required"},
		"password": []string{"required"},
	}

	messages := govalidator.MapData{
		"account":  []string{"required:请输入账号"},
		"password": []string{"required:请输入密码"},
	}

	opts := govalidator.Options{
		Request:         r.Raw(),
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: true,
	}
	v := govalidator.New(opts)
	e := v.Validate()

	if len(e) > 0 {
		return w.Assign("message", e).Jump(r.Raw().URL.Path, e)
	}

	var err error
	_, err = s_user.SignIn(ctx.Request().Session(), r.Form("account", ""), r.Form("password", ""))
	if err != nil {
		return w.Jump(r.Raw().URL.Path, err.Error())
	}

	return w.Redirect(http.StatusFound, "/backend/skeleton")
}

//退出登录
func Out(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	s_user.SignOut(r.Session())
	return w.Redirect(http.StatusFound, "/backend/sign")
}
