package c_user

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	s_user "github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"github.com/gorilla/csrf"
	"net/http"
)

func Forget(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if r.IsMethod("get") {
		return w.
			Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
			Layout("backend/layout/layout.html").
			View(http.StatusOK, "backend/user/forget.html")
	}

	oldPassword := r.Form("oldPassword", "")
	newPassword := r.Form("newPassword", "")
	if len(oldPassword) == 0 {
		return w.Assign("message", "请输入旧密码").Assign("code", 1).Assign("data", "").JSON(http.StatusOK)
	}
	if len(newPassword) == 0 {
		return w.Assign("message", "请输入新密码").Assign("code", 1).Assign("data", "").JSON(http.StatusOK)
	}
	if len(newPassword) < 6 {
		return w.Assign("message", "新密码长度不得小于6个字符").Assign("code", 1).Assign("data", "").JSON(http.StatusOK)
	}

	user := s_user.IsSignIn(r.Session())
	if user == nil {
		return w.Assign("message", "登录信息错误").Assign("code", 1).Assign("data", "").JSON(http.StatusOK)
	}

	if !s_user.ComparePassword(oldPassword, user.Password) {
		return w.Assign("message", "旧密码错误").Assign("code", 1).Assign("data", "").JSON(http.StatusOK)
	}
	var err error
	user.Password, err= s_user.GeneratePassword(newPassword)
	if err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	if _, err := dao.Dao.ID(user.ID).Update(user); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.Assign("message", code.Text(code.SUCCESS)).Assign("code", 0).Assign("data", "").JSON(http.StatusOK)
}
