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
	"strings"
)

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
		w.Assign("folderRegexp", s_attachment.FolderRegexp.String())
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
	return w.Success(result)
}

func CheckMD5(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_attachment.Attachment)
	result.MD5 = r.Param("md5", "")
	if result.MD5 == "" {
		return w.Error(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT, "MD5"))
	}
	if has, err := dao.Dao.Get(result); err != nil {
		return err
	} else if !has {
		return w.Error(code.NOT_FOUND_DATA, code.Text(code.NOT_FOUND_DATA))
	}
	return w.Success(result)
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
		return err
	}
	return w.Success(result)
}

func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	id := r.FormInt("id", 0)
	name := r.Form("name", "")
	if id == 0 || name == "" {
		return w.Error(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT, "id or name"))
	}
	result := new(m_attachment.Attachment)
	result.Name = name
	affected, err := dao.Dao.Id(id).Update(result)
	if err != nil {
		return err
	}
	return w.Success(affected)
}

//单个删除
func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	err := s_attachment.DestroyBatch([]int{r.ParamInt("id", 0)})
	if err != nil {
		return err
	}
	return w.Success()
}

//批量删除
func DestroyBatch(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := r.FormSliceInt("ids")
	err := s_attachment.DestroyBatch(ids)
	if err != nil {
		return err
	}
	return w.Success()
}
