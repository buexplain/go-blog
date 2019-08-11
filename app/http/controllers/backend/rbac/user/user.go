package c_user

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	m_user "github.com/buexplain/go-blog/models/user"
	m_util "github.com/buexplain/go-blog/models/util"
	s_user "github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"github.com/gorilla/csrf"
	"net/http"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	query := m_util.NewQuery("User", ctx).Limit()
	query.Finder.Desc("ID")
	query.WhereKeyword()

	var result m_user.List
	query.Find(&result)

	count := query.Count()

	if query.Error != nil {
		return ctx.Error().WrapServer(query.Error).Location()
	}

	return w.
		Assign("count", count).
		Assign("result", result).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/rbac/user/index.html")
}

func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").View(http.StatusOK, "backend/rbac/user/create.html")
}


func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_user.User)

	if err := r.FormToStruct(result); err != nil {
		return w.JumpBack(err)
	}

	if result.Account == "" {
		return w.JumpBack("请填写账号")
	}

	if !m_util.CheckUnique("User", "Account", result.Account) {
		return w.JumpBack("该账号已存在")
	}

	if result.Password == "" {
		return w.JumpBack("请填写密码")
	}

	if len(result.Password) < 6 {
		return w.JumpBack("密码长度不得小于6个字符")
	}

	if m_user.CheckStatus(result.Status) == false {
		return w.JumpBack("请选择状态")
	}

	if m_user.CheckIdentity(result.Identity) == false {
		return w.JumpBack("请选择身份")
	}

	if result.Nickname == "" {
		result.Nickname = result.Account
	}

	if p, err := s_user.GeneratePassword(result.Password); err != nil {
		return w.JumpBack(err)
	}else {
		result.Password = p
	}

	if _, err := dao.Dao.Insert(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.JumpBack("操作成功")
}


func Edit(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_user.User)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack("参数错误")
	}

	if ok, err := dao.Dao.Get(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	} else if !ok {
		return w.JumpBack("参数错误")
	}

	return w.
		Assign("result", result).
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/rbac/user/create.html")
}


func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_user.User)

	if err := r.FormToStruct(result); err != nil {
		return err
	}

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack("ID错误")
	}

	if result.Account == "" {
		return w.JumpBack("请填写账号")
	}

	if !m_util.CheckUnique("User", "Account", result.Account, result.ID) {
		return w.JumpBack("该账号已存在")
	}

	if result.Password != "" {
		if len(result.Password) < 6 {
			return w.JumpBack("新密码长度不得小于6个字符")
		}
	}

	if m_user.CheckStatus(result.Status) == false {
		return w.JumpBack("请选择状态")
	}

	if m_user.CheckIdentity(result.Identity) == false {
		return w.JumpBack("请选择身份")
	}

	if result.Nickname == "" {
		result.Nickname = result.Account
	}

	if len(result.Password) > 0 {
		if p, err := s_user.GeneratePassword(result.Password); err != nil {
			return w.JumpBack(err)
		}else {
			result.Password = p
		}
	}

	if _, err := dao.Dao.ID(result.ID).Update(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.Jump("/backend/rbac/user", "操作成功")
}


//后台用户修改自己的密码
func Forget(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if r.IsMethod("get") {
		return w.
			Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
			Layout("backend/layout/layout.html").
			View(http.StatusOK, "backend/rbac/user/forget.html")
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
