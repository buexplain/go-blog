package s_services

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/base64"
	"io"
	"strings"
	"xorm.io/xorm"
)

// Import SQL DDL from io.Reader
func ImportDB(dao *xorm.Engine, r io.Reader) ([]sql.Result, error) {
	var results []sql.Result
	var lastError error
	scanner := bufio.NewScanner(r)
	semiColSpliter := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, ';'); i >= 0 {
			return i + 1, data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}

	scanner.Split(semiColSpliter)

	for scanner.Scan() {
		query := strings.Trim(scanner.Text(), " \t\n\r")
		if len(query) > 0 {
			if strings.HasPrefix(query, "INSERT") {
				tmp := strings.LastIndex(query, "(")
				b, err := base64.StdEncoding.DecodeString(query[tmp+1 : len(query)-1])
				if err == nil {
					query = query[0:tmp+1] + string(b) + ")"
				}
			}
			result, err := dao.DB().Exec(query)
			results = append(results, result)
			if err != nil {
				return nil, err
			}
		}
	}
	return results, lastError
}
