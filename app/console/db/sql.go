package db

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

//执行sql语句
var sqlCmd *cobra.Command

func init() {
	sqlCmd = &cobra.Command{
		Use:  "sql",
		Long: "执行sql",
		Run: func(cmd *cobra.Command, args []string) {
			sql := strings.Join(args, " ")
			sql = strings.TrimLeft(sql, " ")
			if len(sql) > 6 && strings.EqualFold(sql[0:6], "select") == true {
				result, err := dao.Dao.Query(sql)
				if err != nil {
					boot.Logger.ErrorF("执行sql失败: %s", err)
					os.Exit(1)
				}
				fmt.Println("********************************** result start **********************************")
				for i, row := range result {
					for k, v := range row {
						fmt.Printf("%-20s %s\n", k, v)
					}
					if i != len(result)-1 {
						fmt.Println("---------------------------------------------------------------------------------")
					}
				}
				fmt.Println("********************************** result end ************************************")
			} else {
				result, err := dao.Dao.Exec(sql)
				if err != nil {
					boot.Logger.ErrorF("执行sql失败: %s", err)
					os.Exit(1)
				}
				lastInsertId, err := result.LastInsertId()
				if err != nil {
					boot.Logger.ErrorF("执行sql失败: %s", err)
					os.Exit(1)
				}
				rowsAffected, err := result.RowsAffected()
				if err != nil {
					boot.Logger.ErrorF("执行sql失败: %s", err)
					os.Exit(1)
				}
				fmt.Printf("lastInsertId %d rowsAffected %d\n", lastInsertId, rowsAffected)
			}
		},
	}
	dbCmd.AddCommand(sqlCmd)
}
