package main

import (
	"github.com/buexplain/go-blog/app/http/boot"
	_ "github.com/buexplain/go-blog/app/http/routers"
)

func main() {
	boot.Run()
}
