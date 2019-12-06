package main

import (
	_ "github.com/buexplain/go-blog/app/http/events"
	_ "github.com/buexplain/go-blog/app/http/routers"
	"time"
)

func init() {
	if location, err := time.LoadLocation("Asia/Shanghai"); err == nil {
		time.Local = location
	} else {
		time.Local = time.FixedZone("CST", 8*3600)
	}
}

func main() {
	h_boot.Run()
}
