package m_user

import (
	m_models "github.com/buexplain/go-blog/models"
)

//用户表
type User struct {
	ID        int           `xorm:"not null pk autoincr INTEGER"`
	CreatedAt m_models.Time `xorm:"DateTime created"`
	UpdatedAt m_models.Time `xorm:"DateTime updated"`
	//账号
	Account string `xorm:"index TEXT"`
	//密码
	Password string `xorm:"TEXT"`
	//昵称
	Nickname string `xorm:"TEXT"`
	//用户身份
	Identity Identity `xorm:"INTEGER"`
	//用户状态
	Status Status `xorm:"INTEGER"`
	//最后登录时间
	LastTime m_models.Time `xorm:"DateTime"`
}

type List []User
