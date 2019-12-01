package console

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd *cobra.Command

func init() {
	RootCmd = &cobra.Command{
		Long: "go-blog framework version 0.1.0",
		Run: func(cmd *cobra.Command, args []string) {
			a_boot.Logger.InfoF("命令行工具，使用帮助: %s --help", os.Args[0])
		},
	}
}
