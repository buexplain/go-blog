package c_oAuth

import (
	"fmt"
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	m_oauth "github.com/buexplain/go-blog/models/oauth"
	m_user "github.com/buexplain/go-blog/models/user"
	s_captcha "github.com/buexplain/go-blog/services/captcha"
	s_oauth "github.com/buexplain/go-blog/services/oauth"
	s_user "github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-slim"
	"github.com/buexplain/go-slim/errors"
	"github.com/gorilla/csrf"
	"net/http"
)

//登录授权回调
func Index(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	//根据请求url的参数获取oauth对象
	oauth, err := s_oauth.New(r)
	if err != nil {
		return err
	}
	//根据请求url的code参数获取token
	var accessResult s_oauth.AccessResult
	origin := s_oauth.OriginURL(r)
	accessResult, err = oauth.GetAccessToken(r)
	if err != nil {
		return w.Jump(origin, err.Error())
	}
	//根据token获取用户信息
	var userInfo s_oauth.UserInfo
	userInfo, err = oauth.GetUserInfo(accessResult.GetAccessToken())
	if err != nil {
		return w.Jump(origin, err.Error())
	}
	if userInfo.GetAccount() == "" {
		return w.Jump(origin, "请授予本站相关权限")
	}
	//检查第三方登录表中是否存在用户信息
	var user m_user.User
	var has bool
	has, err = dao.Dao.Table("Oauth").
		Select("`User`.*").
		Where("`Oauth`.`Account`=?", userInfo.GetAccount()).
		Where("`Oauth`.`Status`=?", oauth.GetStatus()).
		Join("INNER", "User", "User.ID = Oauth.UserID").
		Get(&user)
	if err != nil {
		return w.Jump(origin, err.Error())
	}
	//用户第一次以第三方账号进入本站
	if !has {
		r.Session().Set("oauthUserInfo", userInfo)
		r.Session().Set("oauthUserStatus", int(oauth.GetStatus()))
		w.Assign("account", userInfo.GetAccount())
		w.Assign("nickname", userInfo.GetNickname())
		w.Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw()))
		w.Assign("isSignIn", s_user.IsSignIn(r) != nil)
		return w.View(http.StatusOK, "common/oauth/ask.html")
	}
	//用户再次以第三方账号进入本站
	r.Session().Set(s_user.USER_INFO, user)
	return w.Redirect(http.StatusFound, s_oauth.RedirectURL(r))
}

//注册新用户
func Register(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	account := r.Form("account")
	nickname := r.Form("nickname")
	userInfo, ok := r.Session().Get("oauthUserInfo").(s_oauth.UserInfo)
	if !ok {
		return errors.MarkClient(fmt.Errorf("未找到oauth信息"))
	}
	status := m_oauth.Status(r.Session().GetInt("oauthUserStatus"))
	user, err := s_user.RegisterByOauth(account, nickname, status, userInfo)
	if err != nil {
		return err
	}
	r.Session().Set(s_user.USER_INFO, user)
	r.Session().Del("oauthUserInfo")
	return w.Success(map[string]string{"redirectURL": s_oauth.RedirectURL(r)})
}

//绑定已有账号
func Bind(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	captchaVal := r.Form("captchaVal")
	if captchaVal == "" {
		return code.New(code.INVALID_ARGUMENT, "请输入验证码")
	}
	if s_captcha.Verify(r.Session(), captchaVal) == false {
		return code.New(code.INVALID_ARGUMENT, "验证码错误")
	}
	account := r.Form("account")
	password := r.Form("password")
	userInfo, ok := r.Session().Get("oauthUserInfo").(s_oauth.UserInfo)
	if !ok {
		return errors.MarkClient(fmt.Errorf("未找到oauth信息"))
	}
	status := m_oauth.Status(r.Session().GetInt("oauthUserStatus"))
	user, err := s_user.BindByOauth(account, password, status, userInfo)
	if err != nil {
		return err
	}
	r.Session().Set(s_user.USER_INFO, user)
	r.Session().Del("oauthUserInfo")
	return w.Success(map[string]string{"redirectURL": s_oauth.RedirectURL(r)})
}
