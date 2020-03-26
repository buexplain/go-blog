package helpers

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

//本地时间偏移秒数
var LocalTimeOffsetSeconds string
func init()  {
	//如果更改了服务器时区，则需要重启服务器，以更新改偏移量
	t := time.Now()
	_, o := t.Zone()
	LocalTimeOffsetSeconds = fmt.Sprintf("%+d", o)
}

//字节换算
func FormatSize(size int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	var i int
	f_size := float64(size)
	for i = 0; f_size > 1024 && i < 5; i++ {
		f_size /= 1024
	}
	if i >= 0 && i < len(units) {
		return fmt.Sprintf("%.2f%s", f_size, units[i])
	}
	return ""
}

//递归压缩文件或文件夹
func ZIP(dst *zip.Writer, file string) error {
	queue := []string{}
	queue = append(queue, file)
	for {
		if len(queue) == 0 {
			break
		}
		file := queue[len(queue)-1]
		queue = queue[0 : len(queue)-1]
		fi, err := os.Stat(file)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			fis, err := ioutil.ReadDir(file)
			if err != nil {
				return err
			}
			for _, fi := range fis {
				queue = append(queue, filepath.Join(file, fi.Name()))
			}
			continue
		}
		header, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		header.Name = file
		w, err := dst.CreateHeader(header)
		if err != nil {
			return err
		}
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, f)
		if err != nil {
			if err := f.Close(); err != nil {
				return err
			}
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}

//解压缩
func UnZIP(zipFile, path string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := reader.Close(); err != nil {
			panic(err)
		}
	}()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		filename := filepath.Join(path, file.Name)
		if err := os.MkdirAll(filepath.Dir(filename), 0666); err != nil {
			errRC := rc.Close()
			return fmt.Errorf("%s %s", err, errRC)
		}
		w, err := os.Create(filename)
		if err != nil {
			errRC := rc.Close()
			return fmt.Errorf("%s %s", err, errRC)
		}
		_, err = io.Copy(w, rc)
		errRC := rc.Close()
		errW := w.Close()
		if err != nil {
			return err
		}
		if errRC != nil {
			return errRC
		}
		if errW != nil {
			return errW
		}
	}
	return nil
}

var parseTimeFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02 15",
	"2006-01-02",
	time.RFC3339,
	"2006-01",
	"2006",
}

func ParseInLocation(value string, loc *time.Location) (t time.Time, err error) {
	for _, format := range parseTimeFormats {
		if t, err = time.ParseInLocation(format, value, loc); err == nil {
			break
		}
	}
	return
}

//按本地时间解析，然后返回UTC时间
func ParseTimeLocalToUTC(value string) string {
	for _, format := range parseTimeFormats {
		if t, err := time.ParseInLocation(format, value, time.Local); err == nil {
			return t.UTC().Format(format)
		}
	}
	return ""
}

//按UTC时间解析，然后返回本地时间
func ParseTimeUTCToLocal(value string) string {
	for _, format := range parseTimeFormats {
		if t, err := time.ParseInLocation(format, value, time.UTC); err == nil {
			return t.Local().Format(format)
		}
	}
	return ""
}
