package s_captcha

import (
	"github.com/buexplain/go-fool"
	"github.com/mojocn/base64Captcha"
)

const CaptchaID = "captchaID"

type Option func(*base64Captcha.ConfigDigit)

func SetHeight(height int) Option {
	return func(config *base64Captcha.ConfigDigit) {
		config.Height = height
	}
}

func SetWidth(width int) Option {
	return func(config *base64Captcha.ConfigDigit) {
		config.Width = width
	}
}

func SetMaxSkew(maxSkew float64) Option {
	return func(config *base64Captcha.ConfigDigit) {
		config.MaxSkew = maxSkew
	}
}

func SetDotCount(dotCount int) Option {
	return func(config *base64Captcha.ConfigDigit) {
		config.DotCount = dotCount
	}
}

//设置验证码长度
func SetCaptchaLen(captchaLen int) Option {
	return func(config *base64Captcha.ConfigDigit) {
		config.CaptchaLen = captchaLen
	}
}

//校验验证码
func Verify(session fool.Session, captcha string) bool {
	return base64Captcha.VerifyCaptcha(session.GetString(CaptchaID), captcha)
}

//生成验证码
func Generate(session fool.Session, opts ...Option) (url string) {
	var config = &base64Captcha.ConfigDigit{
		Height:     40,
		Width:      100,
		MaxSkew:    1,
		DotCount:   80,
		CaptchaLen: 4,
	}
	for _, option := range opts {
		option(config)
	}
	id, captchaInstance := base64Captcha.GenerateCaptcha("", *config)
	url = base64Captcha.CaptchaWriteToBase64Encoding(captchaInstance)
	session.Set(CaptchaID, id)
	return
}
