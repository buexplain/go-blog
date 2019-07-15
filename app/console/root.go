package console

import (
	"github.com/buexplain/go-fool/flog"
	"github.com/buexplain/go-fool/flog/formatter"
	"github.com/buexplain/go-fool/flog/handler"
	"github.com/spf13/cobra"
	"os"
)
var Logger *flog.Logger

func init()  {
	Logger = flog.New("std", handler.NewSTD(flog.LEVEL_DEBUG, formatter.NewLine().SetTimeFormat("2006-01-02 15:04:05.99"), flog.LEVEL_WARNING))
}

var RootCmd *cobra.Command

func init() {
	RootCmd = &cobra.Command{
		Long: "go-blog framework version 0.1.0",
		Run: func(cmd *cobra.Command, args []string) {
			Logger.InfoF("命令行工具，使用帮助: %s --help", os.Args[0])
		},
	}
}
