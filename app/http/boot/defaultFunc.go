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
		markerErr := errors.IsMarker(err)
		if markerErr == nil {
			//未知的错误，转为服务端错误，并加上栈信息
			err = fmt.Errorf("%w\n%s", err, debug.Stack())
			markerErr = errors.Mark(err, code.SERVER).(*errors.MrKErr)
		}
		if markerErr.Code() >= code.SERVER {
			defaultServerErrorFunc(ctx, markerErr)
		} else {
			defaultClientErrorFunc(ctx, markerErr)
		}
	} else {
		//未知的恐慌，转为服务端错误，并加上栈信息
		err := fmt.Errorf("%+v\n%s", a, debug.Stack())
		markerErr := errors.MarkServer(err).(*errors.MrKErr)
		defaultServerErrorFunc(ctx, markerErr)
	}
}

//服务端错误处理
func defaultServerErrorFunc(ctx *fool.Ctx, markerErr *errors.MrKErr) {
	ctx.Response().Buffer().Reset()
	isDebug := ctx.App().Debug()
	isJSON := (!ctx.Request().AcceptText() || (ctx.Route() != nil && ctx.Route().HasLabel("json")))
	var responseErr error
	if isJSON {
		//返回json
		if isDebug {
			//返回具体错误
			responseErr = ctx.Response().Error(markerErr.Code(), markerErr.Error(), http.StatusOK)
		} else {
			//屏蔽错误
			responseErr = ctx.Response().Error(markerErr.Code(), http.StatusText(http.StatusInternalServerError), http.StatusOK)
		}
	} else {
		//返回文本
		ctx.Response().Header().Set(constant.HeaderXContentTypeOptions, "nosniff")
		if isDebug {
			responseErr = ctx.Response().Abort(
				http.StatusInternalServerError,
				strings.ReplaceAll(strings.ReplaceAll(markerErr.Error(), "\n", "<br>"), "\t", "&nbsp;&nbsp;&nbsp;&nbsp;"))
		} else {
			responseErr = ctx.Response().Abort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}
	if !isDebug {
		//生产环境，记录错误日志
		Logger.Error(fmt.Sprintf(
			"%s%s %s%s %s%d%s %s",
			"[",
			ctx.Request().Raw().Method,
			ctx.Request().Raw().URL.String(),
			"]",
			"[",
			markerErr.Code(),
			"]",
			markerErr.Error(),
		))
	}
	//响应失败，记录日志
	if responseErr != nil {
		Logger.Error(fmt.Sprintf(
			"%s%s %s%s %s%d%s %s",
			"[",
			ctx.Request().Raw().Method,
			ctx.Request().Raw().URL.String(),
			"]",
			"[",
			markerErr.Code(),
			"]",
			responseErr.Error(),
		))
	}
}

//客户端错误处理
func defaultClientErrorFunc(ctx *fool.Ctx, markerErr *errors.MrKErr) {
	ctx.Response().Buffer().Reset()
	isJSON := (!ctx.Request().AcceptText() || (ctx.Route() != nil && ctx.Route().HasLabel("json")))
	var responseErr error
	if isJSON {
		//返回json
		responseErr = ctx.Response().Error(markerErr.Code(), markerErr.Error(), http.StatusOK)
	} else {
		//返回文本
		ctx.Response().Header().Set(constant.HeaderXContentTypeOptions, "nosniff")
		responseErr = ctx.Response().Abort(http.StatusBadRequest, markerErr.Error())
	}
	//响应失败，记录日志
	if responseErr != nil {
		Logger.Error(fmt.Sprintf(
			"%s%s %s%s %s%d%s %s",
			"[",
			ctx.Request().Raw().Method,
			ctx.Request().Raw().URL.String(),
			"]",
			"[",
			markerErr.Code(),
			"]",
			responseErr.Error(),
		))
	}
}

//错误处理
func defaultErrorFunc(ctx *fool.Ctx, err error) {
	if err == nil {
		return
	}
	markerErr := errors.IsMarker(err)
	if markerErr == nil {
		//未知错误，转为服务端错误
		markerErr = errors.Mark(err, code.SERVER).(*errors.MrKErr)
	}
	if markerErr.Code() >= code.SERVER {
		defaultServerErrorFunc(ctx, markerErr)
	} else {
		defaultClientErrorFunc(ctx, markerErr)
	}
}

//默认路由
func defaultRoute(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ctx.Response().Buffer().Reset()
	ctx.Throw(code.New(code.INVALID_ROUTE))
	return nil
}

type CSRFErrorHandler struct {
}

func (this *CSRFErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isJSON := strings.Contains(r.Header.Get(constant.HeaderAccept), constant.MIMEApplicationJSON)
	if isJSON {
		var content []byte
		content, _ = json.Marshal(map[string]interface{}{"code": code.INVALID_CSRF, "message": code.New(code.INVALID_CSRF), "data": ""})
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
