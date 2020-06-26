package s_backup

import (
	"archive/zip"
	"fmt"
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/helpers"
	s_services "github.com/buexplain/go-blog/services"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"xorm.io/xorm/schemas"
)

var PATH string

func init() {
	PATH = filepath.Join(a_boot.ROOT_PATH, "database/backup")
	if err := os.MkdirAll(PATH, 0666); err != nil {
		log.Fatalln(err)
	}
}

type List []string

func (this List) Has(file string) int {
	for k, v := range this {
		if v == file {
			return k
		}
	}
	return -1
}

func (this List) Size(log string) string {
	path := filepath.Join(PATH, log)
	fi, err := os.Stat(path)
	if err != nil {
		return err.Error()
	}
	return helpers.FormatSize(fi.Size())
}

func GetList() (List, error) {
	tmp, err := filepath.Glob(filepath.Join(PATH, "backup-*秒.zip"))
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(tmp))
	for _, v := range tmp {
		result = append(result, filepath.Base(v))
	}
	sort.Sort(sort.Reverse(sort.StringSlice(result)))
	return List(result), nil
}

//获取上一次备份的时间
func LastBackupTime() *time.Time {
	result, err := GetList()
	if err != nil {
		panic(err)
	}
	if len(result) == 0 {
		return nil
	}
	var t time.Time
	t, err = time.ParseInLocation("backup-2006年01月02日15时04分05秒.zip", result[0], time.Local)
	if err != nil {
		panic(err)
	}
	return &t
}

type Message chan string

func (this Message) Success(message string) {
	message = strings.ReplaceAll(message, "\n", " ")
	message = strings.ReplaceAll(message, "\r", " ")
	this <- fmt.Sprintf("event: success\ndata: %s\n\n", message)
}
func (this Message) Tips(message string) {
	message = strings.ReplaceAll(message, "\n", " ")
	message = strings.ReplaceAll(message, "\r", " ")
	this <- fmt.Sprintf("event: tips\ndata: %s\n\n", message)
}
func (this Message) Fail(err error) {
	s := err.Error()
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	this <- fmt.Sprintf("event: fail\ndata: %s\n\n", s)
}

//备份数据
func Start() <-chan string {
	message := make(Message)
	go func() {
		defer func() {
			close(message)
		}()
		baseName := time.Now().Format("backup-2006年01月02日15时04分05秒") + ".zip"
		backup := filepath.Join(PATH, baseName)
		message.Tips("创建备份文件: " + baseName)
		f, err := os.Create(backup)
		if err != nil {
			message.Fail(err)
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				message.Fail(err)
			}
		}()

		message.Tips("根据备份文件创建zip写入器")
		dst := zip.NewWriter(f)
		defer func() {
			if err := dst.Close(); err != nil {
				message.Fail(err)
				return
			}
		}()
		if a_boot.Config.App.Server.CertFile != "" && a_boot.Config.App.Server.KeyFile != "" {
			message.Tips("往zip写入器中写入证书")
			if err := helpers.ZIP(dst, filepath.Join(a_boot.ROOT_PATH, a_boot.Config.App.Server.CertFile)); err != nil {
				message.Fail(err)
				return
			}
			if err := helpers.ZIP(dst, filepath.Join(a_boot.ROOT_PATH, a_boot.Config.App.Server.KeyFile)); err != nil {
				message.Fail(err)
				return
			}
		}
		message.Tips("往zip写入器中写入附件")
		if err := helpers.ZIP(dst, filepath.Join(a_boot.ROOT_PATH, a_boot.Config.Business.Upload.Save)); err != nil {
			message.Fail(err)
			return
		}

		message.Tips("获取数据库表信息")
		if tables, err := dao.Dao.DBMetas(); err != nil { //获取表信息
			message.Fail(err)
			return
		} else {
			message.Tips("创建数据库导出所需的临时文件")
			tempSqlFile, err := ioutil.TempFile(PATH, "backup-*.sql")
			if err != nil {
				message.Fail(err)
				return
			}
			defer func() {
				if err := os.Remove(tempSqlFile.Name()); err != nil {
					message.Fail(err)
				}
			}()
			message.Tips("导出数据库到临时文件")
			for k, table := range tables {
				message.Tips(fmt.Sprintf("正在导出表: %s 剩余 %d 张表", table.Name, len(tables)-k-1))
				err = s_services.DumpDB(dao.Dao, []*schemas.Table{table}, tempSqlFile, s_services.DUMP_DB_DATA)
				if err != nil {
					message.Fail(err)
					return
				}
			}
			err = tempSqlFile.Close()
			if err != nil {
				message.Fail(err)
				return
			}
			message.Tips("往zip写入器中写入导出的数据库临时文件")
			if err := helpers.ZIP(dst, tempSqlFile.Name()); err != nil {
				message.Fail(err)
				return
			}
			message.Success("备份成功")
		}
	}()
	return message
}
