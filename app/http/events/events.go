package e_events

import (
	h_boot "github.com/buexplain/go-blog/app/http/boot"
	e_syncRbacNode "github.com/buexplain/go-blog/app/http/events/syncRbacNode"
)

func init() {
	h_boot.APP.EventHandler().AddListener(e_syncRbacNode.EVENT_NAME, &e_syncRbacNode.EventListener{})
}