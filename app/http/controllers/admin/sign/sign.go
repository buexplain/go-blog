package c_sign

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/models/user"
	"github.com/buexplain/go-blog/services/captcha"
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"github.com/gorilla/csrf"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"strings"
	"time"
)

//显示登录页面
func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return ctx.Response().Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).View(http.StatusOK, "admin/sign/index.html")
}

//登录
func In(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if s_captcha.Verify(r.Session(), strings.TrimSpace(r.Form("captchaVal", ""))) == false {
		return w.Jump(r.Raw().URL.Path, "验证码错误")
	}

	account := strings.TrimSpace(r.Form("account", ""))
	password := strings.TrimSpace(r.Form("password", ""))

	rules := govalidator.MapData{
		"account": []string{"required"},
		"password": []string{"required"},
	}

	messages := govalidator.MapData{
		"account": []string{"required:请输入账号"},
		"password":    []string{"required:请输入密码"},
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
		return w.Jump(r.Raw().URL.Path, e.Encode())
	}

	user := new(m_user.User)
	user.Account = account
	if has, err := s_user.GetByAccount(user); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}else if !has {
		return w.Jump(r.Raw().URL.Path, "账号或密码错误")
	}

	if !s_user.ComparePassword(password, user.Password) {
		return w.Jump(r.Raw().URL.Path, "账号或密码错误")
	}

	if user.Status != m_user.StatusAllow {
		return w.Abort(http.StatusBadRequest, "账号已被禁用，请联系管理员")
	}

	user.LastTime = time.Now()

	if _, err := s_user.Save(user); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.Redirect(http.StatusFound, "/admin/skeleton")
}
