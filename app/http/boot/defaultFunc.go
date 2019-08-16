package boot

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
	ctx.Response().Buffer().Reset()
	data := fmt.Sprintf("%+v\n%s", a, debug.Stack())
	isDebug := ctx.App().Debug()
	isJSON := (ctx.Request().AcceptJSON() || (ctx.Route() != nil && ctx.Route().HasLabel("json")))
	var err error
	if isJSON {
		ctx.Response().Assign("code", code.SERVER).Assign("message", code.Text(code.SERVER))
		if isDebug {
			ctx.Response().Assign("data", data)
		} else {
			ctx.Response().Assign("data", "")
		}
		err = ctx.Response().JSON(http.StatusOK)
	} else {
		ctx.Response().Header().Set(constant.HeaderXContentTypeOptions, "nosniff")
		if isDebug {
			err = ctx.Response().Plain(http.StatusInternalServerError, data)
		} else {
			err = ctx.Response().Plain(http.StatusInternalServerError, code.Text(code.SERVER))
		}
	}
	if !isDebug {
		ctx.Logger().Error(data, ctx.Request().Raw().Method, ctx.Request().Raw().URL.String())
	}
	if err != nil {
		ctx.Logger().Error(err.Error(), ctx.Request().Raw().Method, ctx.Request().Raw().URL.String())
	}
}

//服务端错误处理
func defaultServerErrorFunc(ctx *fool.Ctx, err error) {
	isDebug := ctx.App().Debug()
	isJSON := (ctx.Request().AcceptJSON() || (ctx.Route() != nil && ctx.Route().HasLabel("json")))
	var responseErr error
	if isJSON {
		ctx.Response().Assign("code", code.SERVER).Assign("message", code.Text(code.SERVER))
		if isDebug {
			ctx.Response().Assign("data", err.Error())
		} else {
			ctx.Response().Assign("data", "")
		}
		responseErr = ctx.Response().JSON(http.StatusOK)
	} else {
		ctx.Response().Header().Set(constant.HeaderXContentTypeOptions, "nosniff")
		if isDebug {
			responseErr = ctx.Response().Abort(http.StatusInternalServerError, err.Error())
		} else {
			responseErr = ctx.Response().Abort(http.StatusInternalServerError, code.Text(code.SERVER))
		}
	}
	if !isDebug {
		ctx.Logger().Error(err.Error(), ctx.Request().Raw().Method, ctx.Request().Raw().URL.String())
	}
	if responseErr != nil {
		ctx.Logger().Error(responseErr.Error(), ctx.Request().Raw().Method, ctx.Request().Raw().URL.String())
	}
}

//客户端错误处理
func defaultClientErrorFunc(ctx *fool.Ctx, err error) {
	isJSON := (ctx.Request().AcceptJSON() || (ctx.Route() != nil && ctx.Route().HasLabel("json")))
	var responseErr error
	if isJSON {
		ctx.Response().Assign("code", code.CLIENT).Assign("message", code.Text(code.CLIENT)).Assign("data", err.Error())
		responseErr = ctx.Response().JSON(http.StatusOK)
	} else {
		ctx.Response().Header().Set(constant.HeaderXContentTypeOptions, "nosniff")
		responseErr = ctx.Response().Abort(http.StatusBadRequest, err.Error())
	}
	if responseErr != nil {
		ctx.Logger().Error(responseErr.Error(), ctx.Request().Raw().Method, ctx.Request().Raw().URL.String())
	}
}

//错误处理
func defaultErrorFunc(ctx *fool.Ctx, err error) {
	if knownErr, ok := err.(*errors.Error); ok && knownErr != nil {
		ctx.Response().Buffer().Reset()
		if knownErr.GetCode() == ctx.Error().ClientCode {
			defaultClientErrorFunc(ctx, err)
		} else {
			defaultServerErrorFunc(ctx, err)
		}
	} else if err != nil {
		ctx.Response().Buffer().Reset()
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
