package c_sysLog

import (
	"github.com/buexplain/go-blog/app/http/boot/code"
	s_sysLog "github.com/buexplain/go-blog/services/sysLog"
	"github.com/buexplain/go-fool"
	"net/http"
	"os"
	"path/filepath"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := s_sysLog.GetList()
	if err != nil {
		return err
	}
	return w.
		Assign("result", result).
		View(http.StatusOK, "backend/server/sysLog/index.html")
}

func Download(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := s_sysLog.GetList()
	if err != nil {
		return err
	}
	if k := result.Has(r.Query("file")); k != -1 {
		return w.Download(filepath.Join(s_sysLog.PATH, result[k]), result[k])
	} else {
		return w.Jump("/backend/sysLog", code.Text(code.INVALID_ARGUMENT))
	}
}

func Show(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := s_sysLog.GetList()
	if err != nil {
		return err
	}
	if k := result.Has(r.Query("file")); k != -1 {
		return w.File(filepath.Join(s_sysLog.PATH, result[k]))
	} else {
		return w.Error(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT))
	}
}

func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := s_sysLog.GetList()
	if err != nil {
		return err
	}
	file := r.Query("file")
	if k := result.Has(file); k != -1 && result[0] != file {
		err := os.Remove(filepath.Join(s_sysLog.PATH, result[k]))
		if err != nil {
			return err
		}
		return w.Redirect(http.StatusFound, "/backend/server/sysLog")
	} else {
		return w.Jump("/backend/server/sysLog", code.Text(code.INVALID_ARGUMENT))
	}
}
