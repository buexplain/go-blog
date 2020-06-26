package staticFile

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-slim"
	"net/http"
	"path/filepath"
	"strings"
)

//静态文件中间件
func Filter(ctx *slim.Ctx, w *slim.Response, r *slim.Request) {
	path := r.Raw().URL.Path
	for _, v := range a_boot.Config.App.StaticFile.Path {
		if strings.HasPrefix(path, v) {
			if a_boot.Config.App.StaticFile.Referer {
				if referer := r.Raw().Header.Get("Referer"); referer != "" && strings.Index(referer, r.Host()) == -1 {
					ctx.Throw(w.Plain(http.StatusForbidden, http.StatusText(http.StatusForbidden)))
				} else {
					ctx.Throw(w.File(filepath.Join(a_boot.ROOT_PATH, path)))
				}
			} else {
				ctx.Throw(w.File(filepath.Join(a_boot.ROOT_PATH, path)))
			}
			return
		}
	}
	ctx.Next()
}
