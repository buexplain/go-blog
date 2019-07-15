package method

import (
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-blog/app/boot"
	"net/http"
	"strings"
)

//方法模拟中间件
func Filter(ctx *fool.Ctx, w *fool.Response, r *fool.Request) {
	if r.Raw().Method == http.MethodDelete {
		if m := r.Query(boot.Config.App.Method.Field); strings.EqualFold(http.MethodGet, m) {
			r.Raw().Method = http.MethodGet
		}
	}else if r.Raw().Method == http.MethodPut || r.Raw().Method == http.MethodPatch {
		if m := r.Form(boot.Config.App.Method.Field); strings.EqualFold(http.MethodPost, m) {
			r.Raw().Method = http.MethodPost
		}
	}
	ctx.Next()
}