package db

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	m_util "github.com/buexplain/go-blog/models/util"
	"github.com/spf13/cobra"
	"os"
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
			boot.Logger.InfoF("开始导入sql文件到数据库: %s", fpath)
			if len(fpath) == 0 {
				boot.Logger.Error("缺失参数: --fpath")
				os.Exit(1)
			}
			_, err := m_util.ImportFromFile(dao.Dao, fpath)
			if err != nil {
				boot.Logger.ErrorF("导入sql文件到数据库失败: %s", err)
				os.Exit(1)
			}
			boot.Logger.Info("导入sql文件到数据库成功")
		},
	}
	importCmd.Flags().StringVarP(&fpath, "fpath", "f", "", "导入到数据库的sql文件")
	dbCmd.AddCommand(importCmd)
}
