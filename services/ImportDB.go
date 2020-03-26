package s_services

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"io"
	"strings"
	"xorm.io/xorm"
)

// Import SQL DDL from io.Reader
func ImportDB(engine *xorm.Engine, r io.Reader) ([]sql.Result, error) {
	var results []sql.Result
	var lastError error
	scanner := bufio.NewScanner(r)

	var inSingleQuote bool
	semiColSpliter := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		for i, b := range data {
			if b == '\'' {
				inSingleQuote = !inSingleQuote
			}
			if !inSingleQuote && b == ';' {
				return i + 1, data[0:i], nil
			}
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
			//解析insert语句的base64编码
			if strings.Index(query, "INSERT INTO") != -1 {
				tmp := strings.LastIndex(query, "VALUES (")
				b, err := base64.StdEncoding.DecodeString(query[tmp+8 : len(query)-1])
				if err == nil {
					query = query[0:tmp+8] + string(b) + ")"
				} else {
					return nil, err
				}
			}
			result, err := engine.Exec(query)
			results = append(results, result)
			if err != nil {
				return nil, err
			}
		}
	}

	return results, lastError
}
