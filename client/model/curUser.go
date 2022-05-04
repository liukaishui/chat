package model

import (
	"chat/common/message"
	"net"
)

type CurUser struct {
	Conn net.Conn
	message.User
}
