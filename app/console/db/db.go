package db

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/console"
	"github.com/spf13/cobra"
)

var dbCmd *cobra.Command = &cobra.Command{
	Use:  "db",
	Long: "数据库操作相关命令",
	Run: func(cmd *cobra.Command, args []string) {
		a_boot.Logger.Info("数据库操作相关命令")
	},
}

func init() {
	console.RootCmd.AddCommand(dbCmd)
}
