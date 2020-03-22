package main

import (
	_ "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot"
	_ "github.com/buexplain/go-blog/app/http/events"
	_ "github.com/buexplain/go-blog/app/http/routers"
)

func main() {
	h_boot.Run()
}
