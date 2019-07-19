package routers

import (
	"github.com/buexplain/go-blog/app/http/boot"
)

func init() {
	regexp(boot.APP.Mux())
	common(boot.APP.Mux())
	admin(boot.APP.Mux())
}
