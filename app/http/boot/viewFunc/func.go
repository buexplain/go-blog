package viewFunc

import (
	"html/template"
	"net/url"
	"strconv"
	"strings"
)

//格式化错误消息
func Message(message interface{}) interface{} {
	if message == nil {
		return template.HTML("")
	}

	//字符串
	if s, ok := message.(string); ok {
		return template.HTML(s)
	}

	//url values
	if values, ok := message.(url.Values); ok {
		var buf strings.Builder
		i := 1
		for _, value := range values {
			for _, v := range value {
				buf.WriteString(strconv.Itoa(i) + "、" + v + "<br>")
				i++
			}
		}
		return template.HTML(buf.String())
	}

	return message
}
