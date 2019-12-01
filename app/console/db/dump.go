package db

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	m_util "github.com/buexplain/go-blog/models/util"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
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
	var mode string
	//保存文件
	var fpath string
	//是否最加写入
	var isAppend bool = false
	dumpCmd = &cobra.Command{
		Use:  "dump",
		Long: "导出数据",
		Run: func(cmd *cobra.Command, args []string) {
			a_boot.Logger.InfoF("开始导出数据库到文件: %s", fpath)
			//获取表信息
			tables, err := dao.Dao.DBMetas()
			if err != nil {
				a_boot.Logger.ErrorF("获取表信息失败: %s", err.Error())
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
				a_boot.Logger.ErrorF("数据库中没有找到表: %s", table)
				os.Exit(1)
			}
			var f *os.File
			if isAppend {
				f, err = os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			}else {
				f, err = os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			}
			if err != nil {
				a_boot.Logger.ErrorF("打开导出到的目标文件失败: %s", err)
				os.Exit(1)
			}
			defer func() {
				_ = f.Close()
			}()
			err = m_util.Dump(dao.Dao, dumpTables, parseMode(mode), f)
			if err != nil {
				a_boot.Logger.ErrorF("导出数据库到文件失败: %s", err)
				os.Exit(1)
			}
			a_boot.Logger.InfoF("导出数据库到文件成功: %s", fpath)
		},
	}

	dumpCmd.Flags().StringVarP(&table, "table", "t", "all", "需要导出的表，默认导出全部表")
	dumpCmd.Flags().StringVarP(&mode,
		"mode",
		"m",
		"1|2|64",
		fmt.Sprintf("需要导出的数据: %d 表结构、%d 表索引、%d 表数据，默认导出全部", m_util.DUMP_STRUCTURE, m_util.DUMP_INDEX, m_util.DUMP_DATA))
	dumpCmd.Flags().StringVarP(&fpath, "fpath", "f", time.Now().Format(filepath.Join(a_boot.ROOT_PATH, "database/2006-01-02-15-04-05.sql")), "导出数据到的文件")
	dumpCmd.Flags().BoolVarP(&isAppend, "append", "a", false, "是否追加写入")
	dbCmd.AddCommand(dumpCmd)
}

func parseMode(m string) int {
	tmp := strings.Split(m, "|")
	result := -1
	for k, v := range tmp {
		i, err := strconv.Atoi(v)
		if err != nil {
			return m_util.DUMP_STRUCTURE | m_util.DUMP_INDEX | m_util.DUMP_DATA
		}
		if k == 0 {
			result = i
			continue
		}
		result = result | i
	}
	return result
}