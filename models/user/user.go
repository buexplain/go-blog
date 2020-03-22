package m_user

import (
	"github.com/buexplain/go-blog/dao"
	"time"
)

//用户表
type User struct {
	ID int `xorm:"not null pk autoincr INTEGER"`
	CreatedAt time.Time `xorm:"DateTime created"`
	UpdatedAt time.Time `xorm:"DateTime updated"`
	//账号
	Account string `xorm:"index TEXT"`
	//密码
	Password string `xorm:"TEXT"`
	//昵称
	Nickname string `xorm:"TEXT"`
	//用户身份
	Identity int `xorm:"INTEGER"`
	//用户状态
	Status int `xorm:"INTEGER"`
	//最后登录时间
	LastTime time.Time `xorm:"DATETIME"`
}

func (this User) StatusText() string {
	return StatusText[this.Status]
}

func (this User) LastTimeText() string {
	return this.LastTime.Format("2006-01-02 15:04:05")
}

type List []User

func GetByAccount(account string) (*User, bool, error) {
	u := new(User)
	has, err := dao.Dao.Where("Account=?", account).Get(u)
	return u, has, err
}

func SaveByID(id int, user *User) (affected int64, err error) {
	affected, err = dao.Dao.ID(id).Update(user)
	return
}
