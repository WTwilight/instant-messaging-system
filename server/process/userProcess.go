package process

import (
	"encoding/json"
	"fmt"
	"github.com/WTwilight/instant-messaging-system/common/message"
	"github.com/WTwilight/instant-messaging-system/server/model"
	"github.com/WTwilight/instant-messaging-system/server/utils"
	"net"
)

type UserProcess struct {
	// 用户的连接
	Conn net.Conn

	// 用户的信息
	User message.User
}

// 处理登录逻辑
func (up *UserProcess)ServerProcessLogin(msg *message.Message) (err error) {
	// 核心代码
	// 1. 先从 msg 中取出 msg.Data ， 并直接反序列化成 LoginMsg
	var loginMsg message.LoginMsg
	err = json.Unmarshal([]byte(msg.Data), &loginMsg)
	if err != nil {
		fmt.Println("json.Unmarshal() 错误:", err)
		return
	}

	// 1. 声明一个 resMsg
	var resMsg message.Message
	resMsg.Type = message.LoginResMsgType

	// 2. 声明一个 loginResMsg, 并完成赋值
	var loginResMsg message.LoginResMsg

	// 对比redis数据库中的数据完成验证
	// 1. 使用 model.UserDaoHandle 到redis 数据库验证
	user, err := model.UserDaoHandle.Login(loginMsg.UserId, loginMsg.UserPwd)
	if err != nil {
		//
		if err == model.ErrorUserNotExist {
			loginResMsg.Code = 500
			loginResMsg.Error = err.Error()
		} else if err == model.ErrorUserWrongPwd {
			loginResMsg.Code = 403
			loginResMsg.Error = err.Error()
		} else {
			loginResMsg.Code = 505
			loginResMsg.Error = "未知错误..."
		}
	} else {	// 登录成功
		fmt.Printf("%s 登录成功!\n", user.UserName)
		up.User = *user
		userMgr.AddOnlineUser(up)

		// 通知其他在线用户,我上线了
		user.UserPwd = ""    // 发送这个用户信息给其他在线用户时，需要将密码去掉
		up.NotifyOthersImOnline(user)

		loginResMsg.OnlineUsers = userMgr.GetAllOnlineUsers()
		loginResMsg.Code = 200
		// 将当前在线用户的userid
		loginResMsg.User = *user
	}

	// 3. 将 loginResMsg 序列化
	data, err := json.Marshal(loginResMsg)
	if err != nil {
		fmt.Println("json.Marshal(loginResMsg) 序列化失败:", err)
		return
	}

	// 4. 将data 赋值给 resMsg
	resMsg.Data = string(data)

	// 5. 对 resMsg 序列化，准备回复给客户端
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal(resMsg) 序列化失败:", err)
		return
	}

	// 6. 发送 data ,将其封装到writePkg 函数中
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	if tf.WritePkg(data) != nil {
		fmt.Println("发送消息失败:", err)
	}
	return
}

// 处理注册逻辑
func (up *UserProcess) ServerProcessRegister(msg *message.Message) (err error) {

	// 1. 将 msg.Data 反序列化成 message.RegisterMsg 结构体
	var registerMsg message.RegisterMsg
	err = json.Unmarshal([]byte(msg.Data), &registerMsg)
	if err != nil {
		fmt.Println("ServerProcessRegister() 中 json.Unmarshal() 错误:", err)
		return
	}

	// 2. 生成用户ID
	registerMsg.User.UserId, err = model.UserDaoHandle.GenerateNewUserId()
	if err != nil {
		fmt.Println("注册用户时生成用户ID错误:", err)
		return
	}

	var resMsg message.Message
	resMsg.Type = message.RegisterMsgType

	var registerResMsg message.RegisterResMsg

	// 访问数据库，完成注册
	err = model.UserDaoHandle.Register(&registerMsg.User)
	if err != nil {
		if err == model.ErrorUserExist {
			fmt.Println("用户已存在")
			registerResMsg.Code = 505
			registerResMsg.Error = model.ErrorUserExist.Error()
		} else {
			registerResMsg.Code = 503
			registerResMsg.Error = "未知错误."
		}
	} else {
		registerResMsg.Code = 200
		registerResMsg.UserId = registerMsg.User.UserId
	}

	data, err := json.Marshal(registerResMsg)
	if err != nil {
		fmt.Println("ServerProcessRegister() 中 json.Marshal(registerResMsg) 错误:", err)
		return
	}

	resMsg.Data = string(data)

	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("ServerProcessRegister() 中 json.Marshal(resMsg) 错误:", err)
		return
	}

	tf := &utils.Transfer{
		Conn: up.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("ServerProcessRegister() 中 tf.WritePkg(data) 错误:", err)
	}
	return
}

// 通知其他在线用户， 我上线了
func (up *UserProcess) NotifyOthersImOnline(user *message.User)  {
	// 遍历 onlineUsers, 逐一发送 NotifyUserStatusMsg
	for _, onlineUser := range userMgr.OnlineUsers {
		if onlineUser.User.UserId == user.UserId {
			continue
		}
		onlineUser.NotifyMeOtherUserOnline(user)
	}
}

// 通知我，其他用户上线了
func (up *UserProcess) NotifyMeOtherUserOnline(user *message.User) {

	var msg message.Message
	msg.Type = message.NotifyUserStatusMsgType

	notifyUserStatusMsg := message.NotifyUserStatusMsg{
		User: *user,
	}
	notifyUserStatusMsg.User.UserStatus= message.UserOnline

	data, err := json.Marshal(notifyUserStatusMsg)
	if err != nil {
		fmt.Println("NotifyMeOtherUserOnline(userId int) 中 json.Marshal(notifyUserStatusMsg) 错误:", err)
		return
	}

	msg.Data = string(data)

	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("NotifyMeOtherUserOnline(userId int) 中 json.Marshal(msg) 错误:", err)
		return
	}

	// 创建一个 Transfer 实例，用于发送消息
	tf := &utils.Transfer{
		Conn: up.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOtherUserOnline(userId int) 中 tf.WritePkg(data) 错误:", err)
		return
	}
}