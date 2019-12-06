package m_util

import (
	"bufio"
	"bytes"
	"database/sql"
	"xorm.io/xorm"
	"io"
	"os"
	"strings"
)

//导入数据
//有此方法是因为 DB.Import() 有 bug
//@link https://xorm.io/xorm/issues/1231#issue-410613530
func Import(dao *xorm.Engine, r io.Reader) ([]sql.Result, error) {
	var results []sql.Result
	var lastError error
	scanner := bufio.NewScanner(r)

	semiColSplit := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.Index(data, []byte(";\n")); i >= 0 {
			return i + 2, data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}

	scanner.Split(semiColSplit)

	for scanner.Scan() {
		query := strings.Trim(scanner.Text(), " \t\n\r")
		if len(query) > 0 {
			result, err := dao.DB().Exec(query)
			results = append(results, result)
			if err != nil {
				return nil, err
			}
		}
	}

	return results, lastError
}

//从文件导入数据
func ImportFromFile(dao *xorm.Engine, fpath string) ([]sql.Result, error) {
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()
	return Import(dao, file)
}
