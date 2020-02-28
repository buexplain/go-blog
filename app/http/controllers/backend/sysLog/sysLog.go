package c_sysLog

import (
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/helpers"
	"github.com/buexplain/go-fool"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var logPath string

func init() {
	logPath = filepath.Join(a_boot.ROOT_PATH, a_boot.Config.Log.Path)
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
	path := filepath.Join(logPath, log)
	fi, err := os.Stat(path)
	if err != nil {
		return err.Error()
	}
	return helpers.FormatSize(fi.Size())
}

func getAllLog() ([]string, error) {
	tmp, err := filepath.Glob(filepath.Join(logPath, "/*.log"))
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(tmp))
	for _, v := range tmp {
		result = append(result, filepath.Base(v))
	}
	sort.Sort(sort.Reverse(stringSlice(result)))
	return result, nil
}

//日志文件排序
type stringSlice []string
func (p stringSlice) Len() int           { return len(p) }
func (p stringSlice) Less(i, j int) bool {
	tmpI := strings.SplitN(p[i], ".", 3)
	tmpJ := strings.SplitN(p[j], ".", 3)
	if len(tmpI) == 2 {
		tmpI = append(tmpI, "-1")
		tmpI[1],tmpI[2] = tmpI[2],tmpI[1]
	}
	if len(tmpJ) == 2 {
		tmpJ = append(tmpJ, "-1")
		tmpJ[1],tmpJ[2] = tmpJ[2],tmpJ[1]
	}
	if len(tmpI) == len(tmpJ) && len(tmpI) == 3 {
		if tmpI[0] == tmpJ[0] {
			i, _ = strconv.Atoi(tmpI[1])
			j, _ = strconv.Atoi(tmpJ[1])
			return i < j
		}
		timeI, _ := time.ParseInLocation("2006-01-02", tmpI[0], time.Local)
		timeJ, _ := time.ParseInLocation("2006-01-02", tmpJ[0], time.Local)
		return timeI.Unix() < timeJ.Unix()
	}
	return p[i] < p[j]
}
func (p stringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	tmp, err := getAllLog()
	if err != nil {
		return err
	}
	result := List(tmp)
	return w.
		Assign("result", result).
		View(http.StatusOK, "backend/sysLog/index.html")
}

func Download(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := getAllLog()
	if err != nil {
		return err
	}
	list := List(result)
	if k := list.Has(r.Query("file")); k != -1 {
		return w.Download(filepath.Join(logPath, list[k]), list[k])
	} else {
		return w.Jump("/backend/sysLog", code.Text(code.INVALID_ARGUMENT))
	}
}

func Show(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := getAllLog()
	if err != nil {
		return err
	}
	list := List(result)
	if k := list.Has(r.Query("file")); k != -1 {
		return w.File(filepath.Join(logPath, list[k]))
	} else {
		return w.Error(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT))
	}
}

func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := getAllLog()
	if err != nil {
		return err
	}
	list := List(result)
	file := r.Query("file")
	if k := list.Has(file); k != -1 && result[0] != file {
		err := os.Remove(filepath.Join(logPath, list[k]))
		if err != nil {
			return err
		}
		return w.Redirect(http.StatusFound, "/backend/sysLog")
	} else {
		return w.Jump("/backend/sysLog", code.Text(code.INVALID_ARGUMENT))
	}
}
