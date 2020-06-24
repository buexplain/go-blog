package db

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	s_services "github.com/buexplain/go-blog/services"
	"github.com/spf13/cobra"
	"os"
)

//从database/init.sql文件中导入数据
var importInitCmd *cobra.Command

func init() {
	importInitCmd = &cobra.Command{
		Use:  "importInit",
		Long: "导入文件 ./database/init.sql 到数据库",
		Run: func(cmd *cobra.Command, args []string) {
			//同步模型到数据库
			if err := syncModels(dao.Dao); err != nil {
				a_boot.Logger.ErrorF("导入 ./database/init.sql 文件到数据库失败: %s", err)
				os.Exit(1)
			}
			fpath := "./database/init.sql"
			a_boot.Logger.Info("开始导入 ./database/init.sql 文件到数据库")
			//打开导出到的目标文件
			file, err := os.Open(fpath)
			if err != nil {
				a_boot.Logger.ErrorF("导入 ./database/init.sql 文件到数据库失败: %s", err)
				os.Exit(1)
			}
			defer func() {
				_ = file.Close()
			}()
			_, err = s_services.ImportDB(dao.Dao, file)
			if err != nil {
				a_boot.Logger.ErrorF("导入 ./database/init.sql 文件到数据库失败: %s", err)
				os.Exit(1)
			}
			if err := dao.Dao.Close(); err != nil {
				a_boot.Logger.ErrorF("导入 ./database/init.sql 文件到数据库失败: %s", err)
				os.Exit(1)
			}
			a_boot.Logger.Info("导入 ./database/init.sql 文件到数据库成功")
		},
	}
	dbCmd.AddCommand(importInitCmd)
}
