package db

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	s_services "github.com/buexplain/go-blog/services"
	"github.com/spf13/cobra"
	"os"
	"xorm.io/core"
)

//从sql文件中导入数据
var importCmd *cobra.Command

func init() {
	//保存文件
	var fpath string
	importCmd = &cobra.Command{
		Use:  "import",
		Long: "导入数据",
		Run: func(cmd *cobra.Command, args []string) {
			a_boot.Logger.InfoF("开始导入sql文件到数据库: %s", fpath)
			if len(fpath) == 0 {
				a_boot.Logger.Error("缺失参数: --fpath")
				os.Exit(1)
			}
			//打开导出到的目标文件
			file, err := os.Open(fpath)
			if err != nil {
				a_boot.Logger.ErrorF("导入sql文件到数据库失败: %s", err)
				os.Exit(1)
			}
			defer func() {
				_ = file.Close()
			}()
			//获取表信息
			var tables []*core.Table
			tables, err = dao.Dao.DBMetas()
			if err != nil {
				a_boot.Logger.ErrorF("导入sql文件到数据库失败: %s", err)
				os.Exit(1)
			}
			//删除所有的表
			for _, table := range tables {
				if _, err := dao.Dao.NewSession().Exec("DROP TABLE "+table.Name); err != nil {
					a_boot.Logger.ErrorF("导入sql文件到数据库失败: %s", err)
					os.Exit(1)
				}
			}
			_,err = s_services.ImportDB(dao.Dao, file)
			if err != nil {
				a_boot.Logger.ErrorF("导入sql文件到数据库失败: %s", err)
				os.Exit(1)
			}
			a_boot.Logger.Info("导入sql文件到数据库成功")
		},
	}
	importCmd.Flags().StringVarP(&fpath, "fpath", "f", "", "导入到数据库的sql文件")
	dbCmd.AddCommand(importCmd)
}

