package process

import "fmt"

var (
	userMgr *UserMgr
)

func GetOnlineUsers() map[int]*UserProcess {
	return userMgr.onlineUsers
}

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// 完成对onlineUsers的添加
func (u *UserMgr) AddOnlineUser(up *UserProcess) {
	u.onlineUsers[up.UserId] = up
}

// 删除
func (u *UserMgr) DelOnlineUser(userId int) {
	delete(u.onlineUsers, userId)
}

// 查看所有用户
func (u *UserMgr) GetAllOnlineUsers() map[int]*UserProcess {
	return u.onlineUsers
}

// 根据id返回用户
func (u *UserMgr) GetOnlineUserById(userId int) (*UserProcess, error) {
	if v, ok := u.onlineUsers[userId]; ok {
		return v, nil
	} else {
		return nil, fmt.Errorf("用户%d 不存在", userId)
	}
}
