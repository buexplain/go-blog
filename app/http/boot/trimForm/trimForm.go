package trimForm

import (
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-fool"
	"strings"
)

//前后空白字符清理
func Filter(ctx *fool.Ctx, w *fool.Response, r *fool.Request)  {
	if err := r.ParseForm(fool.DefaultMaxMemory); err != nil {
		if ctx.Route() != nil && ctx.Route().HasLabel("json") {
			//存在路由，并且路由有json标签，则响应json格式
			ctx.Throw(w.Error(code.INVALID_ARGUMENT, err.Error()))
		} else {
			ctx.Throw(ctx.Response().JumpBack(err.Error()))
		}
	}else {
		trim(r.Raw().Form)
		trim(r.Raw().PostForm)
		if r.Raw().MultipartForm != nil {
			trim(r.Raw().MultipartForm.Value)
		}
		ctx.Next()
	}
}

func trim(data map[string][]string)  {
	if data != nil {
		for key, values := range data {
			for k,v := range values {
				values[k] = strings.TrimSpace(v)
			}
			data[key] = values
		}
	}
}