package c_profile

import (
	"github.com/buexplain/go-fool"
	"net/http"
	"net/http/pprof"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	w.WriteHeader(http.StatusOK)
	name := r.Param("name", "")
	if name == "index" {
		r.Raw().URL.Path = "/debug/profile/"
	}
	switch name {
	case "cmdline":
		pprof.Cmdline(w, r.Raw())
		break
	case "profile":
		pprof.Profile(w, r.Raw())
		break
	case "symbol":
		pprof.Symbol(w, r.Raw())
		break
	case "trace":
		pprof.Trace(w, r.Raw())
		break
	default:
		pprof.Index(w, r.Raw())
	}
	return nil
}
