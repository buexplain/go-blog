package m_attachment

import "github.com/buexplain/go-blog/models"

//文章附件表
type Attachment struct {
	models.Field `xorm:"extends"`
	Name         string `xorm:"TEXT"`
	Path         string `xorm:"TEXT"`
	Ext          string `xorm:"TEXT"`
	Folder       string `xorm:"TEXT"`
	Size         int    `xorm:"INTEGER"`
	MD5          string `xorm:"index TEXT"`
}

type List []Attachment