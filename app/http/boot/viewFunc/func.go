package viewFunc

import (
	"github.com/buexplain/go-validator"
	"html/template"
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

	//校验结果
	if result, ok := message.(*validator.Result); ok {
		return template.HTML(result.ToSimpleString("<br>"))
	}

	return message
}
