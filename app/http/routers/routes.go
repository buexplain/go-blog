package routers

import (
	"github.com/buexplain/go-blog/app/http/boot"
)

func init() {
	regexp(h_boot.APP.Mux())
	common(h_boot.APP.Mux())
	frontend(h_boot.APP.Mux())
	backend(h_boot.APP.Mux())
}
