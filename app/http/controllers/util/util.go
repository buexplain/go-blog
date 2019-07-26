package util

import (
	"net/url"
	"strconv"
	"strings"
)

func UrlValuesString(values url.Values) string {
	if values == nil {
		return ""
	}
	var buf strings.Builder

	for name, value := range values {
		buf.WriteString(name + "：<br>")
		buf.WriteString("：")
		for k, v := range value {
			buf.WriteString(strconv.Itoa(k+1) + "、" + v + "<br>")
		}
	}
	return buf.String()
}
