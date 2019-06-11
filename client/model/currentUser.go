package model

import (
	"github.com/WTwilight/instant-messaging-system/common/message"
	"net"
)

// 客户端很多地方会用到 CurrentUser，因此需要定义一个CurrentUser全局变量
type CurrentUser struct {
	Conn net.Conn
	message.User
}
