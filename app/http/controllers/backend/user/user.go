package c_citizen_user

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/user"
	"github.com/buexplain/go-blog/models/util"
	"github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"github.com/gorilla/csrf"
	"net/http"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	query := m_util.NewQuery("User", ctx).Limit()
	//此处只展示普通用户
	query.Finder.Where("Identity=?", m_user.IdentityCitizen)
	query.Finder.Desc("ID")
	query.Screen()

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
		View(http.StatusOK, "backend/user/index.html")
}

func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").View(http.StatusOK, "backend/user/create.html")
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

	if result.Nickname == "" {
		result.Nickname = result.Account
	}

	if p, err := s_user.GeneratePassword(result.Password); err != nil {
		return w.JumpBack(err)
	}else {
		result.Password = p
	}

	//强制用户身份为普通用户
	result.Identity = m_user.IdentityCitizen

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

	if ok, err := dao.Dao.Where("Identity=?", m_user.IdentityCitizen).Get(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	} else if !ok {
		return w.JumpBack("参数错误")
	}

	return w.
		Assign("result", result).
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/user/create.html")
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

	//强制用户身份为普通用户
	result.Identity = m_user.IdentityCitizen

	if _, err := dao.Dao.ID(result.ID).Where("Identity=?", m_user.IdentityCitizen).Update(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.Jump("/backend/user", "操作成功")
}