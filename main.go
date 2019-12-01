package main

import (
	"github.com/buexplain/go-blog/app/http/boot"
	_ "github.com/buexplain/go-blog/app/http/events"
	_ "github.com/buexplain/go-blog/app/http/routers"
	//_ "golang.org/x/crypto/bcrypt"
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
