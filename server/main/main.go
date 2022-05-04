package main

import (
	"chat/server/model"
	"fmt"
	"net"
)

func process(conn net.Conn) {
	defer conn.Close()

	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务端通讯协程错误，err=", err)
	}
}

func initUserDao() {
	model.MyUserDao = model.NewUserDao(rdb)
}

func main() {
	// 当服务器启动时，就启动redis连接
	initPool("localhost:6379", "", 0)
	initUserDao()

	listen, err := net.Listen("tcp", ":8889")
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
			continue
		}

		go process(conn)
	}
}
