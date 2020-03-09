package c_backup

import (
	"fmt"
	h_boot "github.com/buexplain/go-blog/app/http/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	s_backup "github.com/buexplain/go-blog/services/backup"
	"github.com/buexplain/go-fool"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := s_backup.GetList()
	if err != nil {
		return err
	}
	return w.Assign("result", result).View(http.StatusOK, "backend/backup/index.html")
}

//备份数据
//@link https://developer.mozilla.org/zh-CN/docs/Server-sent_events/EventSource
func Start(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	//接管连接，避免服务端写入超时
	hj, ok := w.Raw().(http.Hijacker)
	if !ok {
		return w.Abort(http.StatusInternalServerError, "doesn't support hijacking")
	}

	//获取具体的连接对象
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		return w.Abort(http.StatusInternalServerError, err.Error())
	}

	//关闭连接
	defer func() {
		if err := conn.Close(); err != nil {
			h_boot.Logger.Error(err.Error())
		}
	}()

	//写入头部
	header := []string{
		"HTTP/1.1 200 OK\r\n",
		"Content-Type: text/event-stream; charset=utf-8\r\n",
		"Cache-Control: no-cache\r\n",
		"Connection: keep-alive\r\n",
		fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC822)),
		"\r\n",
	}
	for _, h := range header {
		if _, err := io.WriteString(bufrw, h); err != nil {
			return nil
		}
	}
	if bufrw.Flush() != nil {
		return nil
	}

	//设置心跳，避免客户端超时
	heart := time.NewTicker(time.Second * 30)
	defer func() {
		heart.Stop()
	}()

	//开始备份
	message := s_backup.Start()
	discard := func() {
		//丢弃所有的消息
		for {
			select {
			case _, ok := <-message:
				if !ok {
					return
				}
			}
		}
	}

	//监听备份进度
	for {
		select {
		case s, ok := <-message:
			if ok {
				if _, err := io.WriteString(bufrw, s); err != nil || bufrw.Flush() != nil {
					//可能是客户端主动关闭了连接，丢弃所有的备份消息
					discard()
					return nil
				}
			} else {
				//任务完成，服务端主动关闭连接
				return nil
			}
		case _ = <-heart.C:
			//写入心跳信息，心跳格式是EventSource文档要求的格式
			if _, err := io.WriteString(bufrw, ": heart\n\n"); err != nil || bufrw.Flush() != nil {
				//写入心跳失败，可能是客户端主动关闭了连接，丢弃所有的备份消息
				discard()
				return nil
			}
		}
	}
}

func Download(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := s_backup.GetList()
	if err != nil {
		return err
	}
	if k := result.Has(r.Query("file")); k != -1 {
		return w.Download(filepath.Join(s_backup.PATH, result[k]), result[k])
	} else {
		return w.Jump("/backend/backup", code.Text(code.INVALID_ARGUMENT))
	}
}

func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := s_backup.GetList()
	if err != nil {
		return err
	}
	file := r.Query("file")
	if k := result.Has(file); k != -1 && result[0] != file {
		err := os.Remove(filepath.Join(s_backup.PATH, result[k]))
		if err != nil {
			return err
		}
		return w.Redirect(http.StatusFound, "/backend/backup")
	} else {
		return w.Jump("/backend/backup", code.Text(code.INVALID_ARGUMENT))
	}
}
