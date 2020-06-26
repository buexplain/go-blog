package code

import (
	"fmt"
	"github.com/buexplain/go-slim/errors"
	"strings"
)

//状态码说明
var text map[int]string = map[int]string{}

const StatusTextUnknown = "未知的状态码"

func Text(code int, moreInfo ...interface{}) string {
	if v, ok := text[code]; ok {
		if len(moreInfo) == 0 {
			return v
		}
		format := v + ":" + strings.Repeat(" %+v", len(moreInfo))
		return fmt.Sprintf(format, moreInfo...)
	} else {
		return StatusTextUnknown
	}
}

func NewM(code int, moreInfo ...interface{}) error {
	if v, ok := text[code]; ok {
		if len(moreInfo) == 0 {
			return errors.Mark(errors.New(v), code)
		}
		format := v + ":" + strings.Repeat(" %+v", len(moreInfo))
		return errors.Mark(fmt.Errorf(format, moreInfo...), code)
	} else {
		return errors.Mark(errors.New(StatusTextUnknown), code)
	}
}

func NewF(code int, format string, a ...interface{}) error {
	return errors.Mark(fmt.Errorf(format, a...), code)
}

func New(code int, message ...string) error {
	var s string
	if len(message) == 0 {
		if v, ok := text[code]; ok {
			s = v
		} else {
			s = StatusTextUnknown
		}
	} else {
		s = message[0]
	}
	return errors.Mark(errors.New(s), code)
}
