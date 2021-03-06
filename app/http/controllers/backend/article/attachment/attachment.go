package c_attachment

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/attachment"
	"github.com/buexplain/go-blog/services"
	"github.com/buexplain/go-blog/services/attachment"
	"github.com/buexplain/go-slim"
	"github.com/gorilla/csrf"
	"net/http"
	"strings"
)

func Index(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
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
		w.Assign("acceptExt", strings.Join(a_boot.Config.Business.Upload.Ext(), ","))
		w.Assign("acceptMimeTypes", strings.Join(a_boot.Config.Business.Upload.MimeType(), ","))
		w.Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw()))
		return w.View(http.StatusOK, "backend/article/attachment/index.html")
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

func CheckMD5(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	result := new(m_attachment.Attachment)
	result.MD5 = r.Param("md5", "")
	if result.MD5 == "" {
		return code.NewM(code.INVALID_ARGUMENT, "MD5")
	}
	if has, err := dao.Dao.Get(result); err != nil {
		return err
	} else if !has {
		return code.New(code.NOT_FOUND_DATA)
	}
	return w.Success(result)
}

func Edit(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	result := new(m_attachment.Attachment)
	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack(code.Text(code.INVALID_ARGUMENT, result.ID))
	}

	if has, err := dao.Dao.Get(result); err != nil {
		return err
	} else if !has {
		return w.JumpBack(code.Text(code.INVALID_ARGUMENT, result.ID))
	}

	if err := result.ReadFile(); err != nil {
		return err
	}
	w.Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw()))
	return w.Assign("result", result).View(http.StatusOK, "backend/article/attachment/create.html")
}

func Update(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	name := r.Form("name")
	content := r.Form("content")
	if name == "" && content == "" {
		return code.NewM(code.INVALID_ARGUMENT, "name or content")
	}

	mod := new(m_attachment.Attachment)
	mod.ID = r.ParamPositiveInt("id")
	if mod.ID <= 0 {
		return code.NewM(code.INVALID_ARGUMENT, "id")
	}
	if has, err := dao.Dao.Get(mod); err != nil {
		return err
	} else if !has {
		return code.NewM(code.INVALID_ARGUMENT, mod.ID)
	}
	if name != "" {
		mod.Name = name
	}

	mod.Content = content

	if err := s_attachment.Update(mod); err != nil {
		return err
	}

	return w.Success()
}

func Upload(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	file, err := r.File("file")
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	var result *m_attachment.Attachment
	result, err = s_attachment.Upload(file, r.Form("folder", ""))
	if err != nil {
		return err
	}
	return w.Success(result)
}

//单个删除
func Destroy(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	err := s_attachment.Destroy([]int{r.ParamInt("id", 0)})
	if err != nil {
		return err
	}
	return w.Success()
}

//批量删除
func DestroyBatch(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	ids := r.FormSliceInt("ids")
	err := s_attachment.Destroy(ids)
	if err != nil {
		return err
	}
	return w.Success()
}

func Download(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	result := new(m_attachment.Attachment)
	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack(code.Text(code.INVALID_ARGUMENT, result.ID))
	}

	if has, err := dao.Dao.Get(result); err != nil {
		return err
	} else if !has {
		return w.JumpBack(code.Text(code.INVALID_ARGUMENT, result.ID))
	}
	result.Ext = "." + result.Ext
	if !strings.HasSuffix(result.Name, result.Ext) {
		result.Name += result.Ext
	}
	return w.Download(result.Path, result.Name)
}
