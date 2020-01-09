package code

import (
	"fmt"
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
		format := v+":"+strings.Repeat(" %+v", len(moreInfo))
		return fmt.Sprintf(format, moreInfo...)
	} else {
		return StatusTextUnknown
	}
}