package s_user

import (
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/user"
	"github.com/buexplain/go-fool"
	"golang.org/x/crypto/bcrypt"
)

func Add(user *m_user.User) (affected int64, err error) {
	return dao.Dao.Insert(user)
}

func GetByAccount(user *m_user.User) (bool, error) {
	return dao.Dao.Where("Account=?", user.Account).Get(user)
}

func GeneratePassword(plaintext string) (password string, err error) {
	var b []byte
	b, err = bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err == nil {
		password = string(b)
	}
	return
}

func ComparePassword(plaintext string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(plaintext)); err != nil {
		return false
	}
	return true
}

func Save(user *m_user.User) (affected int64, err error) {
	affected, err = dao.Dao.ID(user.ID).Update(user)
	return
}

func SignIn(r *fool.Request) {

}