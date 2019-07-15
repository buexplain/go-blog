package main

import (
	"github.com/buexplain/go-blog/app/console"
	_ "github.com/buexplain/go-blog/app/console/asset"
)

func main()  {
	if err := console.RootCmd.Execute(); err != nil {
		console.Logger.ErrorF("console start failed: %s", err)
	}
}