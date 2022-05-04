package process

import (
	"chat/common/message"
	"chat/server/model"
	"chat/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	// 增加一个字段，表示conn是哪个用户的
	UserId int
}

func (u *UserProcess) NotifyOtherOnlineUser(userId int) {
	// 遍历users,发送消息
	for v, user := range userMgr.onlineUsers {
		if v == userId {
			continue
		}
		user.NotifyMeOnline(v)
	}
}

func (u *UserProcess) NotifyMeOnline(userId int) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	// mes.data
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline
	// 将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	mes.Data = string(data)

	// 发送的数据转json
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 发送transfer实例
	tf := &utils.Transfer{
		Conn: u.Conn,
	}
	err = tf.WritePkg(data)
}

func (u *UserProcess) ServerProcessRegister(me *message.Message) (err error) {
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(me.Data), &registerMes)
	if err != nil {
		return
	}

	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		registerResMes.Code = 500
		registerResMes.Error = err.Error()
	} else {
		registerResMes.Code = 200
	}

	b, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	resMes.Data = string(b)

	data, err := json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 发送
	tf := &utils.Transfer{
		Conn: u.Conn,
	}
	err = tf.WritePkg(data)
	return
}

func (u *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	// 先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	// 在声明一个 LoginResMes
	var loginResMes message.LoginResMes

	// 如果用户id=100,密码=123456，认为合法,否则不合法
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}
	} else {
		loginResMes.Code = 200
		// 将登录成功的用户id赋给u
		u.UserId = loginMes.UserId
		userMgr.AddOnlineUser(u)
		// 将当前在线用户的id 放入到UsersId
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println(user, "登录成功")
	}

	b, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	resMes.Data = string(b)

	data, err := json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 发送
	tf := &utils.Transfer{
		Conn: u.Conn,
	}
	err = tf.WritePkg(data)

	// 通知其他用户登陆了
	fmt.Println("通知其他用户登录了")
	u.NotifyOtherOnlineUser(loginMes.UserId)

	return
}
