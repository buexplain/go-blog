package method

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-slim"
	"net/http"
	"strings"
)

//方法欺骗中间件中间件
func Filter(ctx *slim.Ctx, w *slim.Response, r *slim.Request) {
	if r.Raw().Method == http.MethodGet {
		//get 请求切换成 delete 请求
		if m := r.Query(a_boot.Config.App.Method.Field); strings.EqualFold(http.MethodDelete, m) {
			r.Raw().Method = http.MethodDelete
		}
	} else if r.Raw().Method == http.MethodPost {
		m := r.Form(a_boot.Config.App.Method.Field)
		if strings.EqualFold(http.MethodPut, m) {
			//post 请求切换成 put 请求
			r.Raw().Method = http.MethodPut
		} else if strings.EqualFold(http.MethodPatch, m) {
			//post 请求切换成 patch 请求
			r.Raw().Method = http.MethodPatch
		}
	}
	ctx.Next()
}
