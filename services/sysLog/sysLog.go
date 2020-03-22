package s_sysLog

import (
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/helpers"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var PATH string

func init() {
	PATH = filepath.Join(a_boot.ROOT_PATH, a_boot.Config.Log.Path)
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
	tmp, err := filepath.Glob(filepath.Join(PATH, "/*.log"))
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(tmp))
	for _, v := range tmp {
		result = append(result, filepath.Base(v))
	}
	sort.Sort(sort.Reverse(stringSlice(result)))
	return List(result), nil
}

//日志文件排序
type stringSlice []string

func (p stringSlice) Len() int { return len(p) }
func (p stringSlice) Less(i, j int) bool {
	tmpI := strings.SplitN(p[i], ".", 3)
	tmpJ := strings.SplitN(p[j], ".", 3)
	if len(tmpI) == 2 {
		tmpI = append(tmpI, "-1")
		tmpI[1], tmpI[2] = tmpI[2], tmpI[1]
	}
	if len(tmpJ) == 2 {
		tmpJ = append(tmpJ, "-1")
		tmpJ[1], tmpJ[2] = tmpJ[2], tmpJ[1]
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
func (p stringSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
