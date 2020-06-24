package main

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/console"
	_ "github.com/buexplain/go-blog/app/console/asset"
	_ "github.com/buexplain/go-blog/app/console/db"
	"os"
)

func main() {
	defer func() {
		if a := recover(); a != nil {
			a_boot.Logger.ErrorF("%+v", a)
			os.Exit(1)
		}
	}()
	if err := console.RootCmd.Execute(); err != nil {
		a_boot.Logger.ErrorF("console start failed: %s", err)
	}
}
