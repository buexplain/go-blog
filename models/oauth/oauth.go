package m_oauth



import (
	m_models "github.com/buexplain/go-blog/models"
)

//用户三方登录表
type Oauth struct {
	ID        int           `xorm:"not null pk autoincr INTEGER"`
	CreatedAt m_models.Time `xorm:"DateTime created"`
	UpdatedAt m_models.Time `xorm:"DateTime updated"`
	//用户id
	UserID int `xorm:"INTEGER"`
	///第三方的用户的唯一信息
	Account string `xorm:"index TEXT"`
	//第三方的用户的昵称
	Nickname string `xorm:"TEXT"`
	//第三方类型
	Status Status `xorm:"INTEGER"`
}

