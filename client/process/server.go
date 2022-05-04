package process

import (
	"chat/client/utils"
	"chat/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func ShowMenu() {
	fmt.Println("------------恭喜xxx登录成功------------")
	fmt.Println("------------1.显示在线用户列表------------")
	fmt.Println("------------2.发送消息------------")
	fmt.Println("------------3.信息列表------------")
	fmt.Println("------------4.退出系统------------")
	fmt.Printf("请选择(1-4):")

	var key int
	smsProcess := &SmsProcess{}
	_, _ = fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		fmt.Println("显示用户列表")
		outputOnlineUser()
	case 2:
		fmt.Println("发送消息")
		fmt.Printf("请想对大家说点什么:")
		var content string
		fmt.Scanf("%s\n", &content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("退出系统")
		os.Exit(0)
	default:
		fmt.Println("你输入的选项不正确..")
	}
}

func serverProcessMes(conn net.Conn) {
	// 不停的读取服务器消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("正在读取服务端消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("读取服务端消息出错了 err=", err)
			return
		}
		fmt.Println("读到服务器消息,", mes)

		switch mes.Type {
		// 有人上线了
		case message.NotifyUserStatusMesType:
			var notifyUserStatusMes message.NotifyUserStatusMes
			err := json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			if err != nil {
				fmt.Println("json.Unmarshal err=", err)
				continue
			}
			updateUserStatus(&notifyUserStatusMes)
		// 收到消息
		case message.SmsMesType:
			fmt.Println("有人群发消息了")
			outputGroupMes(&mes)
		default:
			fmt.Println("未知消息类型:", mes.Type)
		}
		//fmt.Println(mes)
	}
}
