package process

import (
	utils2 "chat/client/utils"
	"chat/common/message"
	"chat/server/utils"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
}

// 注册
func (u *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	// 1.连接服务器
	conn, err := net.Dial("tcp", ":8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	defer conn.Close()

	// 2.准备通过conn发送消息给服务
	var mes message.Message
	mes.Type = message.RegisterMesType
	// 3.创建一个LoginMes 结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName
	// 4.将loginMes 序列号
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	mes.Data = string(data)

	// 6.将 mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	tf := &utils2.Transfer{
		Conn: conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("tf.WritePkg err=", err)
	}

	mes, err = tf.ReadPkg()
	if err != nil {
		return
	}

	// 将mes.Data 返序列化
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if err != nil {
		fmt.Println(err)
	} else if registerResMes.Code == 200 {
		fmt.Println("注册成功了")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}

	return
}

// 登录
func (u *UserProcess) Login(userId int, userPwd string) (err error) {
	// 1.连接服务器
	conn, err := net.Dial("tcp", ":8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	// 延时关闭
	defer conn.Close()

	// 2.准备通过conn发送消息给服务
	var mes message.Message
	mes.Type = message.LoginMesType
	// 3.创建一个LoginMes 结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	// 4.将loginMes 序列号
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 5.把data赋给mes.Data字段
	mes.Data = string(data)

	// 6.将 mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 把长度通过切片方式发送给服务端
	var pkgLen uint32
	pkgLen = uint32(len(data))
	bytes := make([]byte, 4)
	// [0 0 0 x]
	binary.BigEndian.PutUint32(bytes, pkgLen)

	// 发送长度
	n, err := conn.Write(bytes)
	if n != 4 || err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}

	//fmt.Println("客户端发送长度ok")
	// 发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}

	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		return
	}
	// 将mes.Data 返序列化
	var loginResMes message.LoginResMes
	fmt.Println(mes.Data)
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		fmt.Println("登录失败，json.Unmarshal err=", err)
	} else if loginResMes.Code == 200 {

		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline

		go serverProcessMes(conn)
		// 在线列表
		fmt.Println("在线用户列表如下：")
		for _, v := range loginResMes.UsersId {
			if v == userId {
				continue
			}
			fmt.Println("用户id:\t", v)
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Println("登录成功")
		// 显示我们的登录成功菜单
		for {
			ShowMenu()
		}
	} else {
		fmt.Println("登录失败 err=", loginResMes.Error)
		fmt.Println(loginResMes)
	}
	return
}
