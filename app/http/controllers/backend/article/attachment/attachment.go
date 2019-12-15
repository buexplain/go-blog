package c_attachment

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/attachment"
	"github.com/buexplain/go-blog/models/content"
	"github.com/buexplain/go-blog/services"
	"github.com/buexplain/go-blog/services/attachment"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/upload"
	"github.com/gorilla/csrf"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"xorm.io/xorm"
)

var folderRegexp *regexp.Regexp

func init() {
	folderRegexp = regexp.MustCompile(`^[\w][\w/]*[\w]$`)
}

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if !r.IsAjax() {
		extList, err := s_attachment.GetExtList()
		if err != nil {
			return err
		}
		folderList, err := s_attachment.GetFolderList()
		if err != nil {
			return err
		}
		w.Assign("extList", extList)
		w.Assign("folderList", folderList)
		w.Assign("folderRegexp", folderRegexp.String())
		w.Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw()))
		return w.Layout("backend/layout/layout.html").
			View(http.StatusOK, "backend/article/attachment/index.html")
	}

	query := s_services.NewQuery("Attachment", ctx).Limit().Where()
	ext := r.Query("ext", "all")
	if ext != "all" {
		query.Finder.Where("Ext=?", ext)
	}
	query.Finder.Desc("ID")
	var result m_attachment.List
	var count int64
	query.FindAndCount(&result, &count)
	if query.Error != nil {
		return ctx.Error().WrapServer(query.Error).Location()
	}
	return w.
		Assign("code", code.SUCCESS).
		Assign("message", "操作成功").
		Assign("count", count).
		Assign("data", result).
		JSON(http.StatusOK)
}

func CheckMD5(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_attachment.Attachment)
	result.MD5 = r.Param("md5", "")
	if result.MD5 == "" {
		return w.Assign("data", "").Assign("message", code.Text(code.INVALID_ARGUMENT)).Assign("code", code.INVALID_ARGUMENT).JSON(http.StatusOK)
	}
	if has, err := dao.Dao.Get(result); err != nil {
		return w.Assign("data", "").Assign("message", err.Error()).Assign("code", code.SERVER).JSON(http.StatusOK)
	} else if !has {
		return w.Assign("data", "").Assign("message", code.Text(code.NOT_FOUND_DATA)).Assign("code", code.NOT_FOUND_DATA).JSON(http.StatusOK)
	}
	return w.Assign("data", result).Assign("message", code.Text(code.SUCCESS)).Assign("code", code.SUCCESS).JSON(http.StatusOK)
}

func Upload(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	<- time.After(5*time.Second)
	//获取上传的文件
	file, err := r.File("file")
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	//校验自定义的文件夹名称
	folder := r.Form("folder", "")
	if folder != "" {
		folder = strings.Trim(folder, "/")
		if !folderRegexp.MatchString(folder) || len(folder) > 50 {
			return w.Assign("data", "").Assign("message", "非法文件夹名："+folderRegexp.String()).Assign("code", 1).JSON(http.StatusOK)
		}
		if len(strings.Split(folder, "/")) > 5 {
			return w.Assign("data", "").Assign("message", "文件夹路径过深").Assign("code", 1).JSON(http.StatusOK)
		}
	}

	//得到文件的md5值
	if _, err := file.MD5(); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	//判断文件是否已经上传过了
	result := new(m_attachment.Attachment)
	result.MD5, _ = file.MD5()
	if has, err := dao.Dao.Get(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	} else if has {
		return w.Assign("data", result).Assign("message", code.Text(code.SUCCESS)).Assign("code", code.SUCCESS).JSON(http.StatusOK)
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

	//保存文件
	filePath, err := file.SaveToPath(savePath)
	if err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	//登记文件信息到数据库
	result.Name = file.Name()
	result.Path = filePath
	result.MD5, _ = file.MD5()
	result.Ext = file.Ext()
	if folder == "" {
		result.Folder = "./"
	}else {
		result.Folder = folder
	}
	result.Size = int(file.Size())

	if _, err := dao.Dao.Insert(result); err != nil {
		//插入失败，移除已经保存的文件
		if err := os.Remove(filePath); err != nil {
			return ctx.Error().WrapServer(err).Location()
		}
		return ctx.Error().WrapServer(err).Location()
	}

	//返回结果
	return w.Assign("data", result).Assign("message", code.Text(code.SUCCESS)).Assign("code", code.SUCCESS).JSON(http.StatusOK)
}

func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	id := r.FormInt("id", 0)
	name := r.Form("name", "")
	if id == 0 || name == "" {
		return w.Assign("data", "").Assign("message", code.Text(code.INVALID_ARGUMENT)).Assign("code", code.INVALID_ARGUMENT).JSON(http.StatusOK)
	}
	result := new(m_attachment.Attachment)
	result.Name = name
	affected, err := dao.Dao.Id(id).Update(result)
	if err != nil {
		return err
	}
	return w.Assign("data", affected).Assign("message", code.Text(code.SUCCESS)).Assign("code", code.SUCCESS).JSON(http.StatusOK)
}


//单个删除
func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_content.Content)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack("参数错误")
	}

	if affected, err := dao.Dao.Where("Online=?", m_content.OnlineNo).Delete(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}else if affected > 0 {
		return w.
			Assign("code", code.SUCCESS).
			Assign("message", "操作成功").
			Assign("data", "").
			JSON(http.StatusOK)
	}

	return w.
		Assign("code", 1).
		Assign("message", "操作失败").
		Assign("data", "").
		JSON(http.StatusOK)
}

//批量删除
func DestroyBatch(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := r.FormSliceInt("ids")
	fmt.Println(s_attachment.Destroy(ids))
	return  nil
	if len(ids) == 0 {
		return w.JumpBack("参数错误")
	}

	if affected, err := s_services.DestroyBatch("Content", ids, func(session *xorm.Session) *xorm.Session {
		return  session.Where("Online=?", m_content.OnlineNo)
	}); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}else if affected > 0 {
		return w.
			Assign("code", code.SUCCESS).
			Assign("message", "操作成功").
			Assign("data", "").
			JSON(http.StatusOK)
	}

	return w.
		Assign("code", 1).
		Assign("message", "操作失败").
		Assign("data", "").
		JSON(http.StatusOK)
}