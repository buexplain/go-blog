package e_events

import (
	"github.com/buexplain/go-blog/app/http/boot"
	"github.com/buexplain/go-blog/app/http/events/syncRbacNode"
)

func init() {
	h_boot.APP.EventHandler().AddListener(e_syncRbacNode.EVENT_NAME, &e_syncRbacNode.EventListener{})
}