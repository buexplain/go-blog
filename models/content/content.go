package m_content

import (
	"github.com/buexplain/go-blog/models"
	"html/template"
)

type Content struct {
	ID        int           `xorm:"not null pk autoincr INTEGER"`
	CreatedAt m_models.Time `xorm:"DateTime created"`
	UpdatedAt m_models.Time `xorm:"DateTime updated"`
	//标题
	Title string `xorm:"TEXT"`
	//文章内容，markdown
	Body string `xorm:"TEXT"`
	//
	HTML template.HTML `xorm:"-"`
	//点击量
	Hits int `xorm:"INTEGER"`
	//是否上线
	Online Online `xorm:"INTEGER"`
	//分类id
	CategoryID int `xorm:"INTEGER"`
	//pc端封面图片
	CoverPC string `xorm:"TEXT"`
	//移动端封面图片
	CoverWAP string `xorm:"TEXT"`
	//文章来源
	Origin string `xorm:"TEXT"`
	//页面关键词
	Keywords string `xorm:"TEXT"`
	//页面描述
	Description string `xorm:"TEXT"`
}

type List []*Content
