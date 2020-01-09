package s_attachment

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/attachment"
	"github.com/buexplain/go-fool/errors"
	"github.com/buexplain/go-fool/upload"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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

//上传的时候自定义文件夹的校验正则
var FolderRegexp *regexp.Regexp

func init() {
	FolderRegexp = regexp.MustCompile(`^[\w][\w/]*[\w]$`)
}

//上传文件
func Upload(file *upload.Upload, folder string) (*m_attachment.Attachment, error) {
	if folder != "" {
		folder = strings.Trim(folder, "/")
		if !FolderRegexp.MatchString(folder) {
			return nil, errors.MarkClient(fmt.Errorf("自定义文件夹必须符合规则：%s", FolderRegexp.String()))
		}
		if len(folder) > 50 {
			return nil, errors.MarkClient(fmt.Errorf("自定义文件夹长度必须小于50个字符"))
		}

		if len(strings.Split(folder, "/")) > 5 {
			return nil, errors.MarkClient(fmt.Errorf("自定义文件夹深度不能超过5层"))
		}
	}

	//得到文件的md5值
	if _, err := file.MD5(); err != nil {
		return nil, err
	}

	//判断文件是否已经上传过了
	result := new(m_attachment.Attachment)
	result.MD5, _ = file.MD5()
	if has, err := dao.Dao.Get(result); err != nil {
		return nil, err
	} else if has {
		return result, nil
	}

	//文件没有上传过，设置保存规则
	var savePath string
	if folder == "" {
		//没有自定义文件夹，设置文件名称生成规则，设置保存路径生成规则
		file.SetNameRule(upload.NameRuleRand).SetPathRule(upload.PathRuleDate_2)
		//保存在上传根目录下
		savePath = filepath.Join(a_boot.ROOT_PATH, a_boot.Config.Business.Upload.Save)
	} else {
		//设置保存路径
		savePath = filepath.Join(a_boot.ROOT_PATH, a_boot.Config.Business.Upload.Save, folder)
	}

	file.SetValidateExt(a_boot.Config.Business.Upload.Ext...)

	//保存文件
	if _, err := file.SaveToPath(savePath); err != nil {
		return nil, err
	}

	//登记文件信息到数据库
	result.Name = file.Name()
	result.Path = file.Result()
	result.MD5, _ = file.MD5()
	result.Ext = file.Ext()
	if folder == "" {
		result.Folder = "./"
	} else {
		result.Folder = folder
	}
	result.Size = int(file.Size())

	if _, insertErr := dao.Dao.Insert(result); insertErr != nil {
		//插入失败，移除已经保存的文件
		if removeErr := os.Remove(file.Result()); removeErr != nil {
			return nil, fmt.Errorf("上传附件插入错误: %w 上传附件移除文件错误: %w", insertErr, removeErr)
		}
		return nil, fmt.Errorf("上传附件插入错误: %w", insertErr)
	}

	return result, nil
}

//删除文件
func DestroyBatch(ids []int) error {
	var resultArr m_attachment.List
	err := dao.Dao.Unscoped().In("ID", ids).Find(&resultArr)
	if err != nil {
		return err
	}
	for _, result := range resultArr {
		session := dao.Dao.NewSession()
		defer session.Close()

		if err := session.Begin(); err != nil {
			return err
		}
		if affected, err := session.Unscoped().ID(result.ID).Delete(new(m_attachment.Attachment)); err != nil {
			if err := session.Rollback(); err != nil {
				return err
			}
			return err
		} else if affected > 0 {
			if removeErr := os.Remove(result.Path); removeErr != nil && !os.IsNotExist(removeErr) {
				if rollbackErr := session.Rollback(); rollbackErr != nil {
					return fmt.Errorf("移除附件文件错误: %w，数据库事务回滚错误: %w", removeErr, rollbackErr)
				}
				return removeErr
			}
		}
		if err := session.Commit(); err != nil {
			return err
		}
	}
	return nil
}
