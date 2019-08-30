package routers

import (
	"github.com/buexplain/go-blog/app/http/boot"
)

func init() {
	regexp(boot.APP.Mux())
	common(boot.APP.Mux())
	frontend(boot.APP.Mux())
	backend(boot.APP.Mux())
}
