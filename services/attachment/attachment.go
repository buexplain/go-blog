package s_attachment

import (
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/attachment"
	"os"
	"sort"
)

type ExtList []string
func GetExtList() (result ExtList, err error) {
	err = dao.Dao.Table("Attachment").Distinct("Ext").Select("Ext").OrderBy("ID DESC").Find(&result)
	if err == nil && len(result) > 0 {
		sort.Strings(result)
	}
	return result, err
}

type FolderList []string
func GetFolderList() (result FolderList, err error) {
	err = dao.Dao.Table("Attachment").Distinct("Folder").Select("Folder").OrderBy("ID DESC").Find(&result)
	if err == nil && len(result) > 0 {
		sort.Strings(result)
	}
	return result, err
}

func Destroy(ids []int) error {
	var resultArr m_attachment.List
	err := dao.Dao.Unscoped().In("ID", ids).Find(&resultArr)
	if err != nil {
		return err
	}

	session := dao.Dao.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	for _, result := range resultArr {
		if affected, err := session.Unscoped().ID(result.ID).Delete(new(m_attachment.Attachment)); err != nil {
			if err := session.Rollback(); err != nil {
				return err
			}
			return err
		} else if affected > 0 {
			if err := os.Remove(result.Path); err != nil && !os.IsNotExist(err) {
				if err := session.Rollback(); err != nil {
					return err
				}
				return err
			}
		}
	}

	if err := session.Commit(); err != nil {
		return err
	}
	return nil
}