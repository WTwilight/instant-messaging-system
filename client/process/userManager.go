package process

import (
	"errors"
	"fmt"
	"github.com/WTwilight/instant-messaging-system/client/model"
	"github.com/WTwilight/instant-messaging-system/common/message"
)

var (
	// 声明客户端维护的在线用户map
	onlineUsers = make(map[message.UserIdType]*message.User, 10)

	// 定义一个CurrentUser全局变量, 在用户登录成功后，对其进行初始化
	CurrentUser model.CurrentUser
)

// 处理 NotifyUserStatusMsg
func updateUserStatus(notifyUserStatusMsg *message.NotifyUserStatusMsg) {

	user, ok := onlineUsers[notifyUserStatusMsg.User.UserId]
	if !ok {
		// 原来列表中没有, 则创建一个 message.User
		user = &message.User{
			UserId: 	notifyUserStatusMsg.User.UserId,
			UserName: 	notifyUserStatusMsg.User.UserName,
			UserGender: notifyUserStatusMsg.User.UserGender,
			UserAge: 	notifyUserStatusMsg.User.UserAge,
			UserStatus: notifyUserStatusMsg.User.UserStatus,
		}
	}
	onlineUsers[notifyUserStatusMsg.UserId] = user
	showOnlineUsers()
}

// 在客户端显示当前在线用户
func showOnlineUsers() {

	fmt.Println("当前在线用户:")
	for id, user := range onlineUsers {
		fmt.Printf("    %d[%s]\n", id, user.UserName)
	}
}

func getUserInfoById(userId message.UserIdType) (user *message.User, err error) {
	user, ok := onlineUsers[userId]
	if !ok {
		err = errors.New("用户没有上线")
	}
	return
}