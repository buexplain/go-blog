package s_captcha

import (
	"github.com/buexplain/go-slim"
	"github.com/mojocn/base64Captcha"
	"image/color"
)

//
const CaptchaID = "captchaID"

var store base64Captcha.Store = base64Captcha.DefaultMemStore

//校验验证码
func Verify(session slim.Session, captcha string) bool {
	return store.Verify(session.GetString(CaptchaID), captcha, true)
}

//生成验证码
func Generate(session slim.Session, height int, width int, length int) string {
	driver := base64Captcha.NewDriverString(
		height,
		width,
		1,
		2,
		length,
		"abcdefghjkmnqrstuvxyz123456789",
		&color.RGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: uint8(255)},
		nil,
	)
	c := base64Captcha.NewCaptcha(driver, store)
	id, url, err := c.Generate()
	if err != nil {
		return ""
	}
	session.Set(CaptchaID, id)
	return url
}
