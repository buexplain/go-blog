package db

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	s_services "github.com/buexplain/go-blog/services"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//导出数据
var dumpCmd *cobra.Command

func parseMode(m string) int {
	tmp := strings.Split(m, "|")
	result := -1
	for k, v := range tmp {
		i, err := strconv.Atoi(v)
		if err != nil {
			return s_services.DUMP_DB_STRUCTURE | s_services.DUMP_DB_INDEX | s_services.DUMP_DB_DATA
		}
		if k == 0 {
			result = i
			continue
		}
		result = result | i
	}
	return result
}

func init() {
	//导出模式
	var mode string
	//保存文件
	var fpath string
	dumpCmd = &cobra.Command{
		Use:  "dump",
		Long: "导出数据",
		Run: func(cmd *cobra.Command, args []string) {
			a_boot.Logger.InfoF("开始导出数据库到文件: %s", fpath)
			tables, err := dao.Dao.DBMetas()
			if err != nil {
				a_boot.Logger.ErrorF("获取数据库表信息失败: %s", err)
				os.Exit(1)
			}
			var f *os.File
			f, err = os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				a_boot.Logger.ErrorF("打开导出到的目标文件失败: %s", err)
				os.Exit(1)
			}
			defer func() {
				_ = f.Close()
			}()
			err = s_services.DumpDB(dao.Dao, tables, f, parseMode(mode))
			if err != nil {
				a_boot.Logger.ErrorF("导出数据库到文件失败: %s", err)
				os.Exit(1)
			}
			if err := dao.Dao.Close(); err != nil {
				a_boot.Logger.ErrorF("手动关闭数据库失败: %s", err)
				os.Exit(1)
			}
			a_boot.Logger.InfoF("导出数据库到文件成功: %s", fpath)
		},
	}
	dumpCmd.Flags().StringVarP(&fpath, "fpath", "f", time.Now().Format(filepath.Join(a_boot.ROOT_PATH, "database/2006-01-02-15-04-05.sql")), "导出数据到的文件")
	dumpCmd.Flags().StringVarP(&mode,
		"mode",
		"m",
		fmt.Sprintf("%d|%d|%d", s_services.DUMP_DB_STRUCTURE, s_services.DUMP_DB_INDEX, s_services.DUMP_DB_DATA),
		fmt.Sprintf("需要导出的数据: %d 表结构、%d 表索引、%d 表数据，默认导出全部", s_services.DUMP_DB_STRUCTURE, s_services.DUMP_DB_INDEX, s_services.DUMP_DB_DATA))
	dbCmd.AddCommand(dumpCmd)
}
