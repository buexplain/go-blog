package viewFunc

import (
	"fmt"
	libUrl "net/url"
)

//url生成函数
type url struct {
	data  libUrl.URL
	query libUrl.Values
}

func (this url) String() string {
	if this.query != nil {
		this.data.RawQuery = this.query.Encode()
	}
	return this.data.String()
}

func (this *url) AddParam(key, value string) *url {
	if this.query == nil {
		this.query = this.data.Query()
	}
	this.query.Add(key, value)
	return this
}

func (this *url) SetParam(key string, value interface{}) *url {
	if this.query == nil {
		this.query = this.data.Query()
	}
	this.query.Set(key, fmt.Sprintf("%v", value))
	return this
}

func (this *url) GetParam(key string) string {
	if this.query == nil {
		this.query = this.data.Query()
	}
	return this.query.Get(key)
}

func (this *url) DelParam(key string) *url {
	if this.query == nil {
		this.query = this.data.Query()
	}
	this.query.Del(key)
	return this
}

func (this *url) SetPath(path string) *url {
	this.data.Path = path
	return this
}

func URL(urlObj libUrl.URL) *url {
	return &url{data: urlObj}
}
