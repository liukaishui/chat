package utils

import (
	"chat/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Transfer struct {
	Conn net.Conn
	Buf  []byte
}

func (t *Transfer) ReadPkg() (message.Message, error) {
	buf := make([]byte, 4)
	fmt.Println("读取客户端发送的数据...")
	_, err := t.Conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read err=", err)
		return message.Message{}, err
	}

	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf)

	// 根据pkgLen读取
	buf = make([]byte, pkgLen)
	n, err := t.Conn.Read(buf)
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("conn.Read err=", err)
		return message.Message{}, err
	}

	// 把pkgLen 返序列化->message.Message
	var mes message.Message
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return message.Message{}, err
	}

	return mes, nil
}

func (t *Transfer) WritePkg(data []byte) (err error) {
	// 先发送一个长度
	// 把长度通过切片方式发送给服务端
	var pkgLen uint32
	pkgLen = uint32(len(data))
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, pkgLen)

	// 发送长度
	n, err := t.Conn.Write(bytes)
	if n != 4 || err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}

	//fmt.Println("客户端发送长度ok")
	// 发送消息本身
	_, err = t.Conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}
	return
}
