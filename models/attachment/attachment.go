package m_attachment

import (
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/models"
	"io/ioutil"
	"path/filepath"
)

//文章附件表
type Attachment struct {
	models.Field `xorm:"extends"`
	Name         string `xorm:"TEXT"`
	Path         string `xorm:"TEXT"`
	Ext          string `xorm:"TEXT"`
	Folder       string `xorm:"TEXT"`
	Size         int    `xorm:"INTEGER"`
	MD5          string `xorm:"index TEXT"`
	Content      string `xorm:"-"`
}

type List []Attachment

func (this *Attachment) ReadFile() error {
	b, err := ioutil.ReadFile(filepath.Join(a_boot.ROOT_PATH, this.Path))
	if err != nil {
		return err
	}
	this.Content = string(b)
	return nil
}
