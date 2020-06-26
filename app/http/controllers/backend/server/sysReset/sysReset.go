package c_sysReset

import (
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-slim"
	"github.com/gorilla/csrf"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var PATH string

func init() {
	PATH = filepath.Join(a_boot.ROOT_PATH, "config.toml")
}

func Index(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	config, err := ioutil.ReadFile(PATH)
	if err != nil {
		return err
	}
	w.Assign("config", string(config))
	w.Assign("pid", os.Getpid())
	w.Assign("checkTimeout", a_boot.Config.App.Server.CloseTimedOut.Seconds())
	return w.
		Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		View(http.StatusOK, "backend/server/sysReset/index.html")
}

func Check(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	pid := r.Query("pid")
	if pid == strconv.Itoa(os.Getpid()) {
		return w.Error(1, "正在重启，请稍等...")
	}
	return w.Success(os.Getpid())
}
