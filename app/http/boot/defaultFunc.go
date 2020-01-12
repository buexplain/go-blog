package h_boot

import (
	"encoding/json"
	"fmt"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/constant"
	"github.com/buexplain/go-fool/errors"
	"net/http"
	"runtime/debug"
	"strings"
)

//恐慌恢复
func defaultRecoverFunc(ctx *fool.Ctx, a interface{}) {
	if err, ok := a.(interface{ Error() string }); ok {
		//判断是否实现了错误接口
		if errors.HasMarkerClient(err) {
			//是一个客户端错误
			defaultClientErrorFunc(ctx, err)
		} else {
			//是一个服务端错误
			err = fmt.Errorf("%w\n%s", err, debug.Stack())
			defaultServerErrorFunc(ctx, err)
		}
	} else {
		//未实现错误接口
		err := fmt.Errorf("%+v\n%s", a, debug.Stack())
		defaultServerErrorFunc(ctx, err)
	}
}

//服务端错误处理
func defaultServerErrorFunc(ctx *fool.Ctx, err error) {
	ctx.Response().Buffer().Reset()
	isDebug := ctx.App().Debug()
	isJSON := (ctx.Request().AcceptJSON() || (ctx.Route() != nil && ctx.Route().HasLabel("json")))
	var responseErr error
	if isJSON {
		//返回json
		if isDebug {
			//返回具体错误
			responseErr = ctx.Response().Error(code.SERVER, code.Text(code.SERVER, err), http.StatusInternalServerError)
		} else {
			//屏蔽错误
			responseErr = ctx.Response().Error(code.SERVER, code.Text(code.SERVER), http.StatusInternalServerError)
		}
	} else {
		//返回文本
		ctx.Response().Header().Set(constant.HeaderXContentTypeOptions, "nosniff")
		if isDebug {
			//返回具体错误
			responseErr = ctx.Response().Abort(
				http.StatusInternalServerError,
				strings.ReplaceAll(strings.ReplaceAll(err.Error(), "\n", "<br>"), "\t", "&nbsp;&nbsp;&nbsp;&nbsp;"))
		} else {
			//屏蔽错误
			responseErr = ctx.Response().Abort(http.StatusInternalServerError, code.Text(code.SERVER))
		}
	}
	if !isDebug {
		//生产环境，记录错误日志
		Logger.Error(err.Error(), ctx.Request().Raw().Method, ctx.Request().Raw().URL.String())
	}
	//响应失败，记录日志
	if responseErr != nil {
		Logger.Error(responseErr.Error(), ctx.Request().Raw().Method, ctx.Request().Raw().URL.String())
	}
}

//客户端错误处理
func defaultClientErrorFunc(ctx *fool.Ctx, err error) {
	ctx.Response().Buffer().Reset()
	isJSON := (ctx.Request().AcceptJSON() || (ctx.Route() != nil && ctx.Route().HasLabel("json")))
	var responseErr error
	if isJSON {
		//返回json
		responseErr = ctx.Response().Error(code.CLIENT, code.Text(code.CLIENT, err))
	} else {
		//返回文本
		ctx.Response().Header().Set(constant.HeaderXContentTypeOptions, "nosniff")
		responseErr = ctx.Response().Abort(http.StatusBadRequest, err.Error())
	}
	//响应失败，记录日志
	if responseErr != nil {
		Logger.Error(responseErr.Error(), ctx.Request().Raw().Method, ctx.Request().Raw().URL.String())
	}
}

//错误处理
func defaultErrorFunc(ctx *fool.Ctx, err error) {
	if err == nil {
		return
	}
	if errors.HasMarkerClient(err) {
		//明确的是客户端的错误
		defaultClientErrorFunc(ctx, err)
	} else {
		//通通作为服务端错误
		defaultServerErrorFunc(ctx, err)
	}
}

//默认路由
func defaultRoute(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ctx.Response().Buffer().Reset()
	w.Header().Set(constant.HeaderXContentTypeOptions, "nosniff")
	return w.Plain(http.StatusNotFound, code.Text(code.INVALID_ROUTE))
}

type CSRFErrorHandler struct {
}

func (this *CSRFErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isJSON := strings.Contains(r.Header.Get(constant.HeaderAccept), constant.MIMEApplicationJSON)
	if isJSON {
		var content []byte
		content, _ = json.Marshal(map[string]interface{}{"code": code.INVALID_CSRF, "message": code.Text(code.INVALID_CSRF), "data": ""})
		w.Header().Set(constant.HeaderContentType, constant.MIMEApplicationJSONCharsetUTF8)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(content)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(code.Text(code.INVALID_CSRF)))
	}
}

//默认csrf错误
var defaultCSRFErrorHandler *CSRFErrorHandler

func init() {
	defaultCSRFErrorHandler = new(CSRFErrorHandler)
}
