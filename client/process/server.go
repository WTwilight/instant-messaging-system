package process

import (
	"encoding/json"
	"fmt"
	"github.com/WTwilight/instant-messaging-system/client/utils"
	"github.com/WTwilight/instant-messaging-system/common/message"
	"net"
	"os"
)

// 显示登录成功后的界面
func UserMainMenu() {
	fmt.Println("---xxx的聊天室---")
	fmt.Println("---1. 显示在线用户列表")
	fmt.Println("---2. 发送消息")
	fmt.Println("---3. 信息列表")
	fmt.Println("---4. 退出系统")
	fmt.Println("请选择(1-4):")

	var key int

	fmt.Scanf("%d\n", &key)

	switch key {
	case 1:
		fmt.Println("显示在线用户列表")
		showOnlineUsers()
	case 2:
		sendMsgProcess()
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("退出系统")
		os.Exit(0)

	default:
		fmt.Println("输入有误")
	}
}

func sendMsgProcess() {

	var key int
	var loop = true

	var content string
	sp := SmsProcess{}

	for {
		if !loop {
			break
		}
		fmt.Printf("\n选择发送消息类型:\n")
		fmt.Printf("    1. 群发\n")
		fmt.Printf("    2. 私聊\n")
		fmt.Printf("    3. 返回主菜单\n")
		fmt.Printf("  请选择(1-3):")

		_, _ = fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Printf("群发消息\n")
			fmt.Println("输入要发送消息:")
			_, _ = fmt.Scanf("%s\n", &content)
			_ = sp.SendGroupMsg(content)
		case 2:
			fmt.Printf("发送点对点消息\n")
			sendCtocMsg()
		case key:
			fmt.Printf("返回主菜单")
			loop = false
		default:
			fmt.Printf("输入有误，请重新选择.\n")
		}
	}
	return
}

func sendCtocMsg() {

	// 显示在线用户列表, 供用户选择
	showOnlineUsers()

	// 选择需要发送的用户id
	fmt.Printf("选择给谁发消息[请输入对方ID]:")
	var id message.UserIdType
	_, _ = fmt.Scanf("%d\n", &id)

	// 根据ID发送消息
	sp := &SmsProcess{}
	_ = sp.SendMsgByUserID(id)
}

func serverProcessMsg(conn net.Conn) {

	// 创建一个 Transfer 实例，不停的读取服务器发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}

	for {
		// 客户端等待服务器发送消息
		msg, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("读取消息失败:", err)
			break
		}

		fmt.Println("msg:", msg)
		switch msg.Type {
		case message.NotifyUserStatusMsgType:
			// fmt.Println("有人上线了")
			// 取出 NotifyUserStatusMsg.UserId 和 Status
			var notifyUserStatusMsg message.NotifyUserStatusMsg
			_ = json.Unmarshal([]byte(msg.Data), &notifyUserStatusMsg)
			fmt.Printf("用户  %d[%s] 上线了 \n",
				notifyUserStatusMsg.User.UserId,
				notifyUserStatusMsg.User.UserName)

			// 把用户的信息，状态保存在客户端 map[int]*User
			updateUserStatus(&notifyUserStatusMsg)
		case message.SmsMsgType:
			showGroupMsg(&msg)
		case message.CtocMsgType:
			showCtocMsg(&msg)
		default:
			fmt.Println("未识别的消息类型")
		}
	}
}
