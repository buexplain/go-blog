package code

//状态码说明
var text map[int]string = map[int]string{}

const StatusTextUnknown = "未知的状态码"

func Text(code int) string {
	if v, ok := text[code]; ok {
		return v
	} else {
		return StatusTextUnknown
	}
}