package routers

import (
	"encoding/gob"
	"fmt"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-blog/app/http/boot"
	"os"
	"strconv"
	"time"
)

type Person struct {
	FirstName string
	LastName  string
	Email     string
	Age       int
}

func init() {
	gob.Register(&Person{})
	Register(boot.APP.Mux())
}

func Register(mux *fool.Mux) {
	mux.Regexp("id", `^[1-9][0-9]*$`)
	mux.Any("/", func(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
		var person *Person
		var ok bool

		if person, ok = r.Session().Get("person").(*Person); !ok {
		}
		ctx.Logger().Info("person", person)
		return w.Plain(200, fmt.Sprintf("%+v\n", person))
	})

	mux.Get("/set", func(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
		person := &Person{}
		person.LastName = r.Query("name", "")
		r.Session().Set("person", person)
		return w.Plain(200, "ok")
	})

	mux.Get("/destroy", func(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
		r.Session().Destroy()
		return w.Plain(200, "ok")
	})

	mux.Get("/regenerate", func(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
		r.Session().Regenerate()
		return w.Plain(200, "ok")
	})

	mux.Get("/id", func(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
		return w.Plain(200, r.Session().ID())
	})

	mux.Get("/name", func(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
		return w.Plain(200, r.Session().Name())
	})

	mux.Get("/pid", func(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
		ctx.Logger().Info("test write not found log file")
		return w.Plain(200, strconv.Itoa(os.Getpid()))
	})

	mux.Get("/ppid", func(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
		return w.Plain(200, strconv.Itoa(os.Getppid()))
	})

	mux.Get("/slow", func(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
		ms := r.QueryInt("ms", 10)
		<- time.After(time.Duration(ms)*time.Millisecond)
		return w.Plain(200, strconv.Itoa(os.Getpid()))
	})

	mux.Any("/log", func(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
		msg := r.Query("msg")
		t := time.NewTicker(1*time.Millisecond)
		for {
			<- t.C
			ctx.Logger().Info(msg)
		}

		return w.Plain(200, "ok")
	})
}
