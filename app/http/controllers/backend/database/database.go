package c_database

import (
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	s_backup "github.com/buexplain/go-blog/services/backup"
	s_database "github.com/buexplain/go-blog/services/database"
	"github.com/buexplain/go-slim"
	"github.com/gorilla/csrf"
	"net/http"
	"sort"
	"strings"
	"time"
)

func Index(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	w.Assign("tables", s_database.GetTables())
	w.Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw()))
	return w.View(http.StatusOK, "backend/database/index.html")
}

func SQL(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	sql := r.Form("sql")
	sql = strings.Trim(sql, " ")
	sqlType := r.Form("sqlType", "")
	if sqlType == "" {
		if len(sql) > 6 && strings.EqualFold(sql[0:6], "select") == true {
			sqlType = "query"
		} else {
			sqlType = "exec"
		}
	}

	w.Assign("sql", sql)
	w.Assign("sqlType", sqlType)
	w.Assign("tables", s_database.GetTables())
	w.Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw()))

	lastBackup := s_backup.LastBackupTime()
	if lastBackup == nil || time.Now().Sub(*lastBackup).Hours() > 2 {
		w.Assign("err", "请先备份一次数据")
		goto loop
	}

	if sqlType == "query" {
		b4ExecTime := time.Now()
		rows, err := dao.Dao.Query(sql)
		execDuration := time.Since(b4ExecTime)
		if err != nil {
			w.Assign("err", err)
			goto loop
		}
		w.Assign("execDuration", execDuration)
		if len(rows) > 0 {
			fields := make([]string, 0, len(rows[0]))
			for k, _ := range rows[0] {
				fields = append(fields, k)
			}
			sort.Strings(fields)
			w.Assign("fields", fields)
		}
		w.Assign("rows", rows)
	} else if sqlType == "exec" {
		b4ExecTime := time.Now()
		result, err := dao.Dao.Exec(sql)
		if err != nil {
			w.Assign("err", err)
			goto loop
		}
		execDuration := time.Since(b4ExecTime)
		w.Assign("execDuration", execDuration)
		w.Assign("result", result)
	} else {
		w.Assign("err", code.Text(code.INVALID_ARGUMENT, sqlType))
		goto loop
	}
loop:
	return w.View(http.StatusOK, "backend/database/index.html")
}
