package c_backup

import (
	"archive/zip"
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/helpers"
	s_services "github.com/buexplain/go-blog/services"
	"github.com/buexplain/go-fool"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

var BACKUP_PATH string

func init() {
	BACKUP_PATH = filepath.Join(a_boot.ROOT_PATH, "database/backup")
	if err := os.MkdirAll(BACKUP_PATH, 0666); err != nil {
		panic(err)
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
	path := filepath.Join(BACKUP_PATH, log)
	fi, err := os.Stat(path)
	if err != nil {
		return err.Error()
	}
	return helpers.FormatSize(fi.Size())
}

func getAllBackup() ([]string, error) {
	tmp, err := filepath.Glob(filepath.Join(BACKUP_PATH, "/*.zip"))
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(tmp))
	for _, v := range tmp {
		result = append(result, filepath.Base(v))
	}
	sort.Sort(sort.Reverse(sort.StringSlice(result)))
	return result, nil
}

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	all, err := getAllBackup()
	if err != nil {
		return err
	}
	result := List(all)
	return w.Assign("result", result).View(http.StatusOK, "backend/backup/index.html")
}

func Start(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	//获取当前时间
	t := time.Now().Format("backup-2006年01月02日15时04分05秒")
	//备份数据到具体的文件
	backup := filepath.Join(BACKUP_PATH, t+".zip")
	//创建备份文件
	f, err := os.Create(backup)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	//创建zip写入器
	dst := zip.NewWriter(f)
	defer func() {
		if err := dst.Close(); err != nil {
			panic(err)
		}
	}()

	//往zip写入器中写入附件
	if err := helpers.ZIP(dst, filepath.Join(a_boot.ROOT_PATH, a_boot.Config.Business.Upload.Save)); err != nil {
		return err
	}

	//往zip写入器中写入数据库备份信息
	if tables, err := dao.Dao.DBMetas(); err != nil { //获取表信息
		return err
	} else {
		tempSqlFile, err := ioutil.TempFile(BACKUP_PATH, "backup-*.sql")
		if err != nil {
			return err
		}
		defer func() {
			_ = os.Remove(tempSqlFile.Name())
		}()
		err = s_services.DumpDB(dao.Dao, tables, tempSqlFile)
		if err != nil {
			return err
		}
		err = tempSqlFile.Close()
		if err != nil {
			return err
		}
		//将导出的文件写入到zip写入器中
		if err := helpers.ZIP(dst, tempSqlFile.Name()); err != nil {
			return err
		}
	}

	return w.Success()
}

func Download(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := getAllBackup()
	if err != nil {
		return err
	}
	list := List(result)
	if k := list.Has(r.Query("file")); k != -1 {
		return w.Download(filepath.Join(BACKUP_PATH, list[k]), list[k])
	} else {
		return w.Jump("/backend/backup", code.Text(code.INVALID_ARGUMENT))
	}
}

func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := getAllBackup()
	if err != nil {
		return err
	}
	list := List(result)
	file := r.Query("file")
	if k := list.Has(file); k != -1 && result[0] != file {
		err := os.Remove(filepath.Join(BACKUP_PATH, list[k]))
		if err != nil {
			return err
		}
		return w.Redirect(http.StatusFound, "/backend/backup")
	} else {
		return w.Jump("/backend/backup", code.Text(code.INVALID_ARGUMENT))
	}
}
