package s_user

import (
	"encoding/gob"
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	m_models "github.com/buexplain/go-blog/models"
	m_oauth "github.com/buexplain/go-blog/models/oauth"
	m_user "github.com/buexplain/go-blog/models/user"
	s_oauth "github.com/buexplain/go-blog/services/oauth"
	"github.com/buexplain/go-fool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const USER_INFO = "UserInfo"

func init() {
	gob.Register(&m_user.User{})
}

//生成加密的密码
func GeneratePassword(plaintext string) (password string, err error) {
	var b []byte
	b, err = bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err == nil {
		password = string(b)
	}
	return
}

//校验密码
func ComparePassword(plaintext string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(plaintext)); err != nil {
		return false
	}
	return true
}

//后台管理人员登录
func OfficialSignIn(session fool.Session, account string, password string) (*m_user.User, error) {
	user, has, err := GetByAccount(account)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, code.New(code.INVALID_ARGUMENT, "账号或密码错误")
	}

	if !ComparePassword(password, user.Password) {
		return nil, code.New(code.INVALID_ARGUMENT, "账号或密码错误")
	}

	if user.Status != m_user.StatusAllow {
		return nil, code.New(code.INVALID_AUTH, "账号已被禁用，请联系管理员")
	}

	if user.Identity != m_user.IdentityOfficial {
		return nil, code.New(code.INVALID_AUTH, "非法登录")
	}

	user.LastTime = m_models.Time(time.Now())

	_, err = SaveByID(user.ID, user)

	if err == nil {
		session.Set(USER_INFO, user)
	}

	return user, err
}

//退出登录
func SignOut(session fool.Session) {
	session.Del(USER_INFO)
	session.Destroy()
}

//校验是否登录
func IsSignIn(r *fool.Request) *m_user.User {
	if !r.HasCookie(a_boot.Config.Session.Name) {
		return nil
	}
	user := r.Session().Get(USER_INFO)
	if user == nil {
		return nil
	}
	if u, ok := user.(*m_user.User); ok {
		return u
	}
	return nil
}

func GetByAccount(account string) (*m_user.User, bool, error) {
	u := new(m_user.User)
	has, err := dao.Dao.Where("Account=?", account).Get(u)
	return u, has, err
}

func SaveByID(id int, user *m_user.User) (affected int64, err error) {
	affected, err = dao.Dao.ID(id).Update(user)
	return
}

//根据第三方的账号信息，注册本站用户信息
func RegisterByOauth(account string, nickname string, status m_oauth.Status, userInfo s_oauth.UserInfo) (*m_user.User, error) {
	if account == "" {
		return nil, code.NewM(code.INVALID_ARGUMENT, "account")
	}
	if nickname == "" {
		return nil, code.NewM(code.INVALID_ARGUMENT, "nickname")
	}
	if status.String() == m_models.EnumUNKNOWN {
		return nil, code.NewM(code.INVALID_ARGUMENT, "status")
	}
	var user m_user.User
	//判断是否已经注册过
	has, err := dao.Dao.Table("Oauth").
		Select("`User`.*").
		Where("`Oauth`.`Account`=?", userInfo.GetAccount()).
		Where("`Oauth`.`Status`=?", status).
		Join("INNER", "User", "User.ID = Oauth.UserID").
		Get(&user)
	if err != nil {
		return nil, err
	}
	if has {
		//存在用户，直接返回
		return &user, nil
	}
	//检查需要注册的账号是否唯一
	has, err = dao.Dao.Table("User").
		Where("Account", account).
		Get(&user)
	if has {
		//账号不唯一，直接返回
		return nil, code.New(code.INVALID_ARGUMENT, "账号已存在")
	}
	//生成一个随机的密码
	var password string
	password, err = GeneratePassword(time.Now().String())
	if err != nil {
		return nil, err
	}
	//打开事务
	session := dao.Dao.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return nil, err
	}
	//插入数据到用户表
	user.Status = m_user.StatusAllow
	user.Account = account
	user.Nickname = nickname
	user.Identity = m_user.IdentityCitizen
	user.Password = password
	user.LastTime = m_models.Time(time.Now())
	_, err = dao.Dao.Insert(&user)
	if err != nil {
		if err := session.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	if user.ID == 0 {
		if err := session.Rollback(); err != nil {
			return nil, err
		}
		return nil, code.New(code.SERVER, "插入表 user 失败")
	}
	//插入数据到第三方登录表
	oauth := m_oauth.Oauth{}
	oauth.Nickname = userInfo.GetAccount()
	oauth.Account = userInfo.GetAccount()
	oauth.UserID = user.ID
	oauth.Status = status
	_, err = dao.Dao.Insert(&oauth)
	if err != nil {
		if err := session.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	if oauth.ID == 0 {
		if err := session.Rollback(); err != nil {
			return nil, err
		}
		return nil, code.New(code.SERVER, "插入表 oauth 失败")
	}
	//提交事务
	if err := session.Commit(); err != nil {
		return nil, err
	}
	return &user, nil
}

//将用户的第三方账号信息与本站用户信息进行绑定
func BindByOauth(account string, password string, status m_oauth.Status, userInfo s_oauth.UserInfo) (*m_user.User, error) {
	if account == "" {
		return nil, code.NewM(code.INVALID_ARGUMENT, "account")
	}
	if password == "" {
		return nil, code.NewM(code.INVALID_ARGUMENT, "password")
	}
	if status.String() == m_models.EnumUNKNOWN {
		return nil, code.NewM(code.INVALID_ARGUMENT, "status")
	}
	var oauth m_oauth.Oauth
	//判断是否已经绑定过
	has, err := dao.Dao.Table("Oauth").
		Where("`Oauth`.`Account`=?", userInfo.GetAccount()).
		Where("`Oauth`.`Status`=?", status).
		Get(&oauth)
	if err != nil {
		return nil, err
	}
	if has {
		return nil, code.NewM(code.CLIENT, "账号已被绑定到 "+userInfo.GetNickname())
	}
	var user *m_user.User
	user, has, err = GetByAccount(account)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, code.New(code.CLIENT, "账号或密码错误")
	}
	if !ComparePassword(password, user.Password) {
		return nil, code.New(code.CLIENT, "账号或密码错误")
	}
	if user.Status != m_user.StatusAllow {
		return nil, code.New(code.CLIENT, "账号已被禁用，请联系管理员")
	}
	oauth.Nickname = userInfo.GetAccount()
	oauth.Account = userInfo.GetAccount()
	oauth.UserID = user.ID
	oauth.Status = status
	_, err = dao.Dao.Insert(&oauth)
	if err != nil {
		return nil, err
	}
	if oauth.ID == 0 {
		return nil, code.New(code.SERVER, "插入表 oauth 失败")
	}
	return user, nil
}
