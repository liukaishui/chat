package main

import (
	"chat/client/process"
	"fmt"
	"os"
)

// 定义两个变量，一个表示用户id，一个表示用户密码
var (
	userId   int
	userPwd  string
	userName string
)

func main() {
	// 接收用户的选择
	var key int

	for {
		fmt.Println("--------欢迎登录多人聊天系统--------")
		fmt.Println("1 登录聊天系统")
		fmt.Println("2 注册用户")
		fmt.Println("3 退出系统")
		fmt.Print("请选择(1-3):")

		_, _ = fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Print("请输入用户的id:")
			_, _ = fmt.Scanf("%d\n", &userId)
			fmt.Print("请输入用户的密码:")
			_, _ = fmt.Scanf("%s\n", &userPwd)
			up := &process.UserProcess{}
			up.Login(userId, userPwd)
		case 2:
			fmt.Println("注册用户")
			fmt.Printf("请输入用户id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Printf("请输入用户密码:")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Printf("请输入用户昵称:")
			fmt.Scanf("%s\n", &userName)
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Printf("\n您的输入有误，请重新输入\n")
		}
	}
}
