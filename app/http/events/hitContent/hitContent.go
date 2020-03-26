package e_hitContent

import (
	h_boot "github.com/buexplain/go-blog/app/http/boot"
	"github.com/buexplain/go-blog/dao"
	m_content "github.com/buexplain/go-blog/models/content"
	"github.com/buexplain/go-event"
)

//事件名称
const EVENT_NAME = "hitContent"

//事件监听者
type EventListener struct {
}

func (this *EventListener) Handle(e *event.Event) {
	contentID, ok := e.Data.(int)
	if !ok || contentID == 0 {
		return
	}
	_, err := dao.Dao.Where("ID=?", contentID).Incr("Hits", 1).Update(new(m_content.Content))
	if err != nil {
		h_boot.Logger.ErrorF("文章浏览量累加失败: %s", err)
	}
}
