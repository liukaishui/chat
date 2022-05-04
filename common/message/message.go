package message

const (
	LoginMesType            = "loginMes" // 登录
	LoginResMesType         = "LoginResRes"
	RegisterMesType         = "RegisterMes" // 注册
	RegisterResMesType      = "RegisterResMesType"
	NotifyUserStatusMesType = "NotifyUserStatusMes" // 通知
	SmsMesType              = "SmsMes"
)

// 定义用户状态常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type LoginMes struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

type LoginResMes struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	UsersId []int
}

type RegisterMes struct {
	User User `json:"user"`
}

type RegisterResMes struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type NotifyUserStatusMes struct {
	UserId int `json:"userId"`
	Status int `json:"status"`
}

// 增加一个SmsMes 发送消息
type SmsMes struct {
	Content string `json:"content"`
	User
}
