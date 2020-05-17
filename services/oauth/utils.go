package s_oauth

import (
	"math/rand"
	"net/url"
)

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomString(n int, allowedChars ...[]rune) string {
	var letters []rune

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func setQueryString(uri string, key string, value string) string {
	if uri == "" {
		return ""
	}
	tmp, err := url.Parse(uri)
	if err != nil {
		return uri
	}
	query := tmp.Query()
	query.Set(key, value)
	tmp.RawQuery = query.Encode()
	return tmp.String()
}