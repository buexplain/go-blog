package db

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/dao/util"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"xorm.io/core"
)

//导出数据
var dumpCmd *cobra.Command

func init() {
	//导出的表
	var table string
	//导出模式
	var mold int
	//保存文件
	var fpath string
	dumpCmd = &cobra.Command{
		Use:  "dump",
		Long: "导出数据",
		Run: func(cmd *cobra.Command, args []string) {
			boot.Logger.InfoF("开始导出数据库到文件: %s", fpath)
			//获取表信息
			tables, err := dao.Dao.DBMetas()
			if err != nil {
				log.Panicln(err)
				boot.Logger.ErrorF("获取表信息失败: %s", err.Error())
				os.Exit(1)
			}
			var dumpTables []*core.Table
			if !strings.EqualFold(table, "all") {
				dumpTables = make([]*core.Table, 0, len(tables))
				for _, v := range tables {
					if strings.EqualFold(v.Name, table) {
						dumpTables = append(dumpTables, v)
					}
				}
			} else {
				dumpTables = tables
			}
			if len(dumpTables) == 0 {
				boot.Logger.ErrorF("数据库中没有找到表: %s", table)
				os.Exit(1)
			}
			var f *os.File
			f, err = os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				boot.Logger.ErrorF("打开导出到的目标文件失败: %s", err)
				os.Exit(1)
			}
			defer func() {
				_ = f.Close()
			}()
			err = util.Dump(dao.Dao, dumpTables, mold, f)
			if err != nil {
				boot.Logger.ErrorF("导出数据库到文件失败: %s", err)
				os.Exit(1)
			}
			boot.Logger.InfoF("导出数据库到文件成功: %s", fpath)
		},
	}
	dumpCmd.Flags().StringVarP(&table, "table", "t", "all", "需要导出的表，默认导出全部表")
	dumpCmd.Flags().IntVarP(&mold,
		"mold",
		"m",
		util.DUMP_STRUCTURE|util.DUMP_INDEX|util.DUMP_DATA,
		fmt.Sprintf("需要导出的数据: %d 表结构、%d 表索引、%d 表数据，默认导出全部: 1|2|64", util.DUMP_STRUCTURE, util.DUMP_INDEX, util.DUMP_DATA))
	dumpCmd.Flags().StringVarP(&fpath, "fpath", "f", time.Now().Format(filepath.Join(boot.ROOT_PATH, "database/2006-01-02-15-04-05.sql")), "导出数据到的文件")
	dbCmd.AddCommand(dumpCmd)
}
