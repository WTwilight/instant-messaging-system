package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/WTwilight/instant-messaging-system/client/utils"
	"github.com/WTwilight/instant-messaging-system/common/message"
	"net"
)

type UserProcess struct {
	//
}

func (up *UserProcess)Login(userId message.UserIdType, userPwd string) (err error) {

	// 1. 连接到服务器端
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial 错误:", err)
		return
	}

	defer func() {
		if conn.Close() != nil {
			fmt.Println("conn.Close() 关闭连接失败!")
		}
	}()

	// 2. 准备通过 conn 发送消息给服务器
	var msg message.Message
	msg.Type = message.LoginMsgType

	// 3. 创建一个 LoginMsg 结构体
	var loginMsg message.LoginMsg
	loginMsg.UserId = userId
	loginMsg.UserPwd = userPwd

	// 4. 将 loginMsg 结构体序列化
	loginMsgString, err := json.Marshal(loginMsg)
	if err != nil {
		fmt.Println("json.Marshal 错误:", err)
		return
	}

	// 5. 把data 赋给 msg.Data 字段
	msg.Data = string(loginMsgString)

	// 6. 将 msg 序列化
	msgString, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Marshal 错误:", err)
		return
	}

	// 7. 此时 data 是要发送的消息
	// 7.1 先把 data 的长度发送给服务器
	// 先获取 msgString 的长度，再转成一个表示长度的 []byte 切片
	var msgStringLen uint32
	msgStringLen = uint32(len(msgString))
	fmt.Printf("msgStringLen 的值为 %d\n", msgStringLen)
	lenBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(lenBytes[:4], msgStringLen)
	fmt.Println("lenBytes:", lenBytes[:4])

	// 发送长度
	n, err := conn.Write(lenBytes[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(lenBytes) 失败:", err)
		return
	}
	fmt.Println("客户端发送消息长度成功:",len(msgString), "内容:", string(msgString))

	// 发送消息本身
	_, err = conn.Write(msgString)
	if err != nil {
		fmt.Println("conn.Write(msgString) 错误:", err)
		return
	}

	// 这里还需要处理服务器返回消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	msg, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) 错误:", err)
		return
	}

	var loginResMsg message.LoginResMsg
	err = json.Unmarshal([]byte(msg.Data), &loginResMsg)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(msg.Data), &loginResMsg) 反序列化失败:", err)
		return
	}
	if loginResMsg.Code == 200 {
		fmt.Println("登录成功")
		// 初始化 CurrentUser 变量
		CurrentUser.Conn = conn
		CurrentUser.User = loginResMsg.User
		CurrentUser.User.UserStatus = message.UserOnline
		fmt.Println("我的账户信息:", CurrentUser.User)

		// 显示当前在线用户列表
		fmt.Println("当前在线用户列表:")
		for _, u := range loginResMsg.OnlineUsers {
			if u.UserId == userId {	// 在线用户列表不显示自己
				continue
			}
			fmt.Printf(" %d[%s]\n", u.UserId, u.UserName)
			// 将在线用户信息保存到 onlineUsers 中
			onlineUsers[u.UserId] = &u
		}
		fmt.Printf("\n\n")

		go serverProcessMsg(conn)

		for  {
			UserMainMenu()
		}
	} else {
		fmt.Println("登录失败:", loginResMsg.Error)
	}
	return
}

func (up *UserProcess) Register(user *message.User) (err error) {

	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial 错误:", err)
		return
	}

	defer func() {
		if conn.Close() != nil {
			fmt.Println("conn.Close() 关闭连接失败!")
		}
	}()

	var msg message.Message
	msg.Type = message.RegisterMsgType

	// 3. 创建一个 LoginMsg 结构体
	var registerMsg message.RegisterMsg
	registerMsg.User = *user

	// 4. 将 registerMsg 结构体序列化
	registerMsgString, err := json.Marshal(registerMsg)
	if err != nil {
		fmt.Println("json.Marshal 错误:", err)
		return
	}
	msg.Data = string(registerMsgString)

	// 6. 将 msg 序列化
	msgString, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Marshal 错误:", err)
		return
	}

	tf := &utils.Transfer{
		Conn: conn,
	}

	err = tf.WritePkg(msgString)
	if err != nil {
		fmt.Println("注册用户 tf.WritePkg() 错误:", err)
		return
	}

	msg, err = tf.ReadPkg()    // 此处接收的 msg 就是 RegisterResMsg
	if err != nil {
		fmt.Println("注册用户 tf.ReadPkg() 错误:", err)
		return
	}

	var registerResMsg message.RegisterResMsg
	err = json.Unmarshal([]byte(msg.Data), &registerResMsg)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(msg.Data), &loginResMsg) 反序列化失败:", err)
		return
	}
	if registerResMsg.Code == 200 {
		fmt.Printf("注册成功, 你的id为 %d 请登录.\n\n", registerResMsg.UserId)

	} else {
		fmt.Println("注册失败:", registerResMsg.Error)
	}
	return
}