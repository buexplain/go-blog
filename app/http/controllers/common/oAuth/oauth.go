package c_oAuth

import (
	"fmt"
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	m_oauth "github.com/buexplain/go-blog/models/oauth"
	m_user "github.com/buexplain/go-blog/models/user"
	s_captcha "github.com/buexplain/go-blog/services/captcha"
	s_oauth "github.com/buexplain/go-blog/services/oauth"
	s_user "github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
	"github.com/gorilla/csrf"
	"net/http"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	oauth, err := s_oauth.New(r)
	if err != nil {
		return err
	}
	var accessResult s_oauth.AccessResult
	accessResult, err = oauth.GetAccessToken(r)
	if err != nil {
		return err
	}
	var userInfo s_oauth.UserInfo
	userInfo, err = oauth.GetUserInfo(accessResult.GetAccessToken())
	if err != nil {
		return err
	}
	if userInfo.GetAccount() == "" {
		return errors.MarkClient(fmt.Errorf("请授予本站相关权限"))
	}
	var user m_user.User
	var has bool
	has, err = dao.Dao.Table("Oauth").
		Select("`User`.*").
		Where("`Oauth`.`Account`=?", userInfo.GetAccount()).
		Where("`Oauth`.`Status`=?", oauth.GetStatus()).
		Join("INNER", "User", "User.ID = Oauth.UserID").
		Get(&user)
	if err != nil {
		return err
	}
	//用户第一次以第三方账号进入本站
	if !has {
		r.Session().Set("oauth_user_info", userInfo)
		r.Session().Set("oauth_user_status", int(oauth.GetStatus()))
		w.Assign("account", userInfo.GetAccount())
		w.Assign("nickname", userInfo.GetNickname())
		w.Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw()))
		return w.View(http.StatusOK, "common/oauth/ask.html")
	}
	//用户再次以第三方账号进入本站
	oauth_after_url := r.Session().GetString("oauth_after_url")
	if oauth_after_url == "" {
		oauth_after_url = "/"
	}
	r.Session().Set(s_user.SESSION_ID, user)
	return w.Redirect(http.StatusFound, oauth_after_url)
}

func Register(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	account := r.Form("account")
	nickname := r.Form("nickname")
	userInfo, ok := r.Session().Get("oauth_user_info").(s_oauth.UserInfo)
	if !ok {
		return fmt.Errorf("未找到oauth信息")
	}
	status := m_oauth.Status(r.Session().GetInt("oauth_user_status"))
	user, err := s_user.RegisterByOauth(account, nickname, status, userInfo)
	if err != nil {
		return err
	}
	r.Session().Set(s_user.SESSION_ID, user)
	oauth_after_url := r.Session().GetString("oauth_after_url")
	if oauth_after_url == "" {
		oauth_after_url = "/"
	}
	return w.Success(map[string]string{"oauth_after_url":oauth_after_url})
}

func Bind(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	captchaVal := r.Form("captchaVal")
	if captchaVal == "" {
		return errors.MarkClient(fmt.Errorf("请输入验证码"))
	}
	if s_captcha.Verify(r.Session(), captchaVal) == false {
		return errors.MarkClient(fmt.Errorf("验证码错误"))
	}
	account := r.Form("account")
	password := r.Form("password")
	userInfo, ok := r.Session().Get("oauth_user_info").(s_oauth.UserInfo)
	if !ok {
		return fmt.Errorf("未找到oauth信息")
	}
	status := m_oauth.Status(r.Session().GetInt("oauth_user_status"))
	user, err := s_user.BindByOauth(account, password, status, userInfo)
	if err != nil {
		return err
	}
	r.Session().Set(s_user.SESSION_ID, user)
	oauth_after_url := r.Session().GetString("oauth_after_url")
	if oauth_after_url == "" {
		oauth_after_url = "/"
	}
	return w.Success(map[string]string{"oauth_after_url":oauth_after_url})
}