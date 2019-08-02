package db

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/menu"
	"github.com/buexplain/go-blog/models/user"
	"github.com/spf13/cobra"
)

//同步模型到数据库
var syncCmd *cobra.Command

func init() {
	syncCmd = &cobra.Command{
		Use:  "sync",
		Long: "同步models到数据库",
		Run: func(cmd *cobra.Command, args []string) {
			boot.Logger.Info("开始同步models到数据库")
			err := dao.Dao.Sync2(
				new(m_user.User),
				new(m_menu.Menu),
			)
			if err != nil {
				boot.Logger.ErrorF("同步models到数据库失败: %s", err)
			} else {
				boot.Logger.Info("同步models到数据库成功")
			}
		},
	}
	dbCmd.AddCommand(syncCmd)
}
