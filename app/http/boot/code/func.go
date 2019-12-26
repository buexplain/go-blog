package code

import (
	"github.com/buexplain/go-fool"
	"net/http"
)

func Success(ctx *fool.Ctx, data ...interface{}) error {
	if !ctx.Response().Store().Has("data") {
		if len(data) == 0 {
			data = append(data, "")
		}
	}
	if len(data) > 0 {
		ctx.Response().Assign("data", data[0])
	}
	return ctx.Response().Assign("message", Text(SUCCESS)).
		Assign("code", SUCCESS).
		JSON(http.StatusOK)
}

func Error(ctx *fool.Ctx, code int, message interface{}, data ...interface{}) error {
	if !ctx.Response().Store().Has("data") {
		if len(data) == 0 {
			data = append(data, "")
		}
	}
	if len(data) > 0 {
		ctx.Response().Assign("data", data[0])
	}
	if message == nil {
		message = Text(code)
	}
	return ctx.Response().
		Assign("code", code).
		Assign("message", message).
		JSON(http.StatusOK)
}