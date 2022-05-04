package process

import (
	"chat/common/message"
	"encoding/json"
	"fmt"
)

func outputGroupMes(mes *message.Message) {
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	// 显示信息
	fmt.Printf("用户ID:\t%d 对大家说:\t%s\n", smsMes.UserId, smsMes.Content)
}
