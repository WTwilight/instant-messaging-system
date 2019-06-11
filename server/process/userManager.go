package process

import (
	"fmt"
	"github.com/WTwilight/instant-messaging-system/common/message"
)

// 因为UserMgr 实例在服务器中有且只有一个
// 且在很多地方都要用的，因此将其 声明为一个全局变量

var userMgr *UserMgr

type UserMgr struct {
	OnlineUsers map[message.UserIdType]*UserProcess
}

// 初始化userMgr
func init() {
	userMgr = &UserMgr{
		OnlineUsers: make(map[message.UserIdType]*UserProcess, 1024),
	}
}

// 向 userMgr 中添加元素
func (um *UserMgr) AddOnlineUser(up *UserProcess) {
	um.OnlineUsers[up.User.UserId] = up
}

// 从 userMgr 中删除元素
func (um *UserMgr) DelOfflineUser(userId message.UserIdType) {
	delete(um.OnlineUsers, userId)
}

// 返回当前所有在线用户
func (um *UserMgr) GetAllOnlineUser() map[message.UserIdType]*UserProcess {
	return um.OnlineUsers
}

func (um *UserMgr) GetAllOnlineUsers() (onlineUsers []message.User) {

	for _, up := range um.OnlineUsers {
		onlineUsers = append(onlineUsers, up.User)
	}
	return
}

func (um *UserMgr) GetOnlineUserById(userId message.UserIdType) (up *UserProcess, err error) {
	up, ok := um.OnlineUsers[userId]
	if !ok {
		fmt.Println("查找的用户不在线.")
		err = fmt.Errorf("用户 %d 不在线", userId)
	}
	return
}