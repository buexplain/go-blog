package e_events

import (
	"github.com/buexplain/go-blog/app/http/boot"
	e_hitContent "github.com/buexplain/go-blog/app/http/events/hitContent"
	"github.com/buexplain/go-blog/app/http/events/syncRbacNode"
)

func init() {
	h_boot.Bus.AddListener(e_syncRbacNode.EVENT_NAME, &e_syncRbacNode.EventListener{})
	h_boot.Bus.AddListener(e_hitContent.EVENT_NAME, &e_hitContent.EventListener{})
}
