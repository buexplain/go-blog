package main

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/console"
	_ "github.com/buexplain/go-blog/app/console/asset"
	_ "github.com/buexplain/go-blog/app/console/db"
	"os"
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
	defer func() {
		if a := recover(); a != nil {
			a_boot.Logger.ErrorF("console run failed: %s", a)
			os.Exit(1)
		}
	}()
	if err := console.RootCmd.Execute(); err != nil {
		a_boot.Logger.ErrorF("console start failed: %s", err)
	}
}
