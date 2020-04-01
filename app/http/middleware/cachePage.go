package middleware

import (
	h_boot "github.com/buexplain/go-blog/app/http/boot"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/constant"
	"io"
	"net/http"
)

//缓存页面
func CachePage(ctx *fool.Ctx, w *fool.Response, r *fool.Request) {
	//开发者模式，不做缓存
	if ctx.App().Debug() {
		ctx.Next()
		return
	}
	//不对ajax请求进行缓存，不对非get请求进行缓存
	if r.IsAjax() || !r.IsMethod(http.MethodGet) {
		ctx.Next()
		return
	}
	//该中间件必须是路由级别的中间件
	if ctx.Route() == nil {
		h_boot.Logger.Warning("routing level Middleware: CachePage")
		ctx.Next()
		return
	}
	//如果该路由要求返回json，则不做缓存
	if ctx.Route().HasLabel("json") {
		ctx.Next()
		return
	}
	key := r.Raw().URL.Path
	if h_boot.Cache.Exists(key) {
		reader, writer, err := h_boot.Cache.Get(key)
		if err != nil {
			//缓存异常，移除缓存
			_ = h_boot.Cache.Remove(key)
			h_boot.Logger.ErrorF("read cache error: %s", err)
			//进入下一个中间件
			ctx.Next()
			return
		}
		defer func() {
			if writer != nil {
				_ = writer.Close()
			}
			_ = reader.Close()
		}()
		//读取缓存
		w.Header().Set(constant.HeaderContentType, constant.MIMETextHTMLCharsetUTF8)
		w.WriteHeader(http.StatusOK)
		_, err = io.Copy(w, reader)
		if err != nil {
			ctx.Throw(err)
		}
	}else {
		//缓存不存在，进入下一个中间件
		ctx.Next()
		//http返回200，写入缓存
		if w.StatusCode() == http.StatusOK {
			reader, writer, err := h_boot.Cache.Get(key)
			if err != nil {
				//缓存异常，移除缓存
				_ = h_boot.Cache.Remove(key)
				h_boot.Logger.ErrorF("open cache error: %s", err)
				return
			}
			defer func() {
				if writer != nil {
					_ = writer.Close()
				}
				_ = reader.Close()
			}()
			//写入缓存
			b := w.Buffer().Bytes()
			_, err = writer.Write(b)
			if err != nil {
				ctx.Throw(err)
			}
		}
	}
}