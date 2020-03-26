package s_user

import (
	"encoding/gob"
	"github.com/buexplain/go-blog/dao"
	m_models "github.com/buexplain/go-blog/models"
	m_user "github.com/buexplain/go-blog/models/user"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const SESSION_ID = "SESSION_ID"

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
		return nil, errors.MarkClient(errors.New("账号或密码错误"))
	}

	if !ComparePassword(password, user.Password) {
		return nil, errors.MarkClient(errors.New("账号或密码错误"))
	}

	if user.Status != m_user.StatusAllow {
		return nil, errors.MarkClient(errors.New("账号已被禁用，请联系管理员"))
	}

	if user.Identity != m_user.IdentityOfficial {
		return nil, errors.MarkClient(errors.New("非法登录"))
	}

	user.LastTime = m_models.Time(time.Now())

	_, err = SaveByID(user.ID, user)

	if err == nil {
		session.Set(SESSION_ID, user)
	}

	return user, err
}

//退出登录
func SignOut(session fool.Session) {
	session.Del(SESSION_ID)
	session.Destroy()
}

//校验是否登录
func IsSignIn(session fool.Session) *m_user.User {
	user := session.Get(SESSION_ID)
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
