package c_attachment

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/attachment"
	"github.com/buexplain/go-blog/services"
	"github.com/buexplain/go-blog/services/attachment"
	"github.com/buexplain/go-fool"
	"github.com/gorilla/csrf"
	"net/http"
	"regexp"
	"strings"
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
		w.Assign("acceptExt", strings.Join(a_boot.Config.Business.Upload.Ext, ","))
		w.Assign("acceptMimeTypes", strings.Join(a_boot.Config.Business.Upload.MimeTypes, ","))
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
		return query.Error
	}
	w.Assign("count", count)
	return code.Success(ctx, result)
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
	file, err := r.File("file")
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	result, err := s_attachment.Upload(file, r.Form("folder", ""))
	if err != nil {
		return w.Assign("data", "").Assign("message", err.Error()).Assign("code", code.SERVER).JSON(http.StatusOK)
	}
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
	err := s_attachment.DestroyBatch([]int{r.ParamInt("id", 0)})
	if err != nil {
		return w.Assign("data", "").Assign("message", err.Error()).Assign("code", code.SERVER).JSON(http.StatusOK)
	}
	return w.
		Assign("code", code.SUCCESS).
		Assign("message", "操作成功").
		Assign("data", "").
		JSON(http.StatusOK)
}

//批量删除
func DestroyBatch(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := r.FormSliceInt("ids")
	err := s_attachment.DestroyBatch(ids)
	if err != nil {
		return w.Assign("data", "").Assign("message", err.Error()).Assign("code", code.SERVER).JSON(http.StatusOK)
	}
	return w.
		Assign("code", code.SUCCESS).
		Assign("message", "操作成功").
		Assign("data", "").
		JSON(http.StatusOK)

}