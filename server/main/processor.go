package main

import (
	"chat/common/message"
	process2 "chat/server/process"
	"chat/server/utils"
	"errors"
	"fmt"
	"net"
)

type Processor struct {
	Conn net.Conn
}

// 编写一个serverProcessMes 函数
// 功能: 根据客户端发送消息类型不同，决定调用哪个函数来处理
func (p *Processor) ServerProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	// 处理登录逻辑
	case message.LoginMesType:
		up := &process2.UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessLogin(mes)
	// 处理注册逻辑
	case message.RegisterMesType:
		up := &process2.UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessRegister(mes)
	// 发送消息
	case message.SmsMesType:
		smsProcess := &SmsProcess{}
		smsProcess.SendGroupMes(mes)
	default:
		err = errors.New("消息类型不存在，无法处理")
	}
	return
}

func (p *Processor) process2() (err error) {
	tf := &utils.Transfer{
		Conn: p.Conn,
	}

	for {
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("readPkg err=", err)
			return err
		}
		err = p.ServerProcessMes(&mes)
		if err != nil {
			fmt.Println("serverProcessMes err=", err)
			return err
		}
	}
}
