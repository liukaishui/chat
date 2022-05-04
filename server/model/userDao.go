package model

import (
	"chat/common/message"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

// 服务器启动的时候，就创建一个userDao
var (
	MyUserDao *UserDao
)

// 完成对user结构体的各种操作
type UserDao struct {
	Redis *redis.Client
}

// 工厂模式
func NewUserDao(rdb *redis.Client) *UserDao {
	return &UserDao{
		Redis: rdb,
	}
}

func (u *UserDao) Register(user *message.User) (err error) {
	fmt.Println(user)
	_, err = u.GetUserById(user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		return
	}
	ctx := context.Background()
	err = u.Redis.HSet(ctx, "users", strconv.Itoa(user.UserId), data).Err()
	if err != nil {
		return
	}
	return
}

func (u *UserDao) GetUserById(id int) (*User, error) {
	ctx := context.Background()
	res, err := u.Redis.HGet(ctx, "users", strconv.Itoa(id)).Result()
	if err != nil {
		err = ERROR_USER_NOTEXISTS
		return &User{}, err
	}

	// 把查询到的res，反序列化成实例
	user := &User{}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		return user, err
	}

	return user, err
}

func (u *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	user, err = u.GetUserById(userId)
	if err != nil {
		return
	}
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}
