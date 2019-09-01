package c_util

import (
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-fool"
	"net/http"
)

//响应成功的json
func Success(w *fool.Response, data ...interface{}) error {
	w.Assign("message", code.Text(code.SUCCESS)).Assign("code", code.SUCCESS)
	if len(data) > 0 {
		w.Assign("data", data[0])
	}else {
		w.Assign("data", "")
	}
	return w.JSON(http.StatusOK)
}

//响应错误的json
func Error(w *fool.Response, message string, data ...interface{}) error {
	w.Assign("message", message).Assign("code", 1)
	if len(data) > 0 {
		w.Assign("data", data[0])
	}else {
		w.Assign("data", "")
	}
	return w.JSON(http.StatusOK)
}