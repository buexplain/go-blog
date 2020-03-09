package helpers

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

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
