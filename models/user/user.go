package m_user

import (
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models"
	"time"
)

type User struct {
	models.Field `xorm:"extends"`
	//账号
	Account string `xorm:"index TEXT"`
	//密码
	Password string `xorm:"TEXT"`
	//昵称
	Nickname string `xorm:"TEXT"`
	//用户状态
	Status int `xorm:"INTEGER"`
	//最后登录时间
	LastTime time.Time `xorm:"DATETIME"`
}

const (
	StatusAllow = iota + 1
	StatusDeny
)

var StatusText = map[int]string{
	StatusAllow: "允许",
	StatusDeny:  "禁止",
}

func GetByAccount(account string) (*User, bool, error) {
	u := new(User)
	has, err := dao.Dao.Where("Account=?", account).Get(u)
	return u, has, err
}

func SaveByID(id int, user *User) (affected int64, err error) {
	affected, err = dao.Dao.ID(id).Update(user)
	return
}
