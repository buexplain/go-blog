package m_content

import (
	"github.com/buexplain/go-blog/models"
	"html/template"
)

type Content struct {
	models.Field `xorm:"extends"`
	Title        string        `xorm:"TEXT"`
	Body         string        `xorm:"TEXT"`
	HTML         template.HTML `xorm:"TEXT"`
	Hits         int           `xorm:"INTEGER"`
	Online       int           `xorm:"INTEGER"`
	CategoryID     int           `xorm:"INTEGER"`
	CoverPC      string        `xorm:"TEXT"`
	CoverWAP     string        `xorm:"TEXT"`
	Origin       string        `xorm:"TEXT"`
}

func (this Content) OnlineText() string {
	return OnlineText[this.Online]
}

type List []*Content

const (
	OnlineYes = iota + 1
	OnlineNo
)

var OnlineText = map[int]string{
	OnlineYes: "已上线",
	OnlineNo:  "已下线",
}
