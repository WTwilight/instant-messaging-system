package main

import (
	"fmt"
	"github.com/WTwilight/instant-messaging-system/common/message"
	"github.com/WTwilight/instant-messaging-system/server/process"
	"github.com/WTwilight/instant-messaging-system/server/utils"
	"net"
)

// 创建 Processor结构体
type Processor struct {
	Conn net.Conn
}

// 编写一个 serverProcessMsg() 函数
// 根据客户端发送消息的类型，决定调用哪个函数来处理
func (p *Processor)serverProcessMsg(msg *message.Message) (err error) {

	switch msg.Type {
	case message.LoginMsgType:
		//调用处理登录的逻辑
		up := &process.UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessLogin(msg)
		if err != nil {
			fmt.Println("serverProcessLogin() 登录失败:", err)
			return
		}
	case message.RegisterMsgType:
		// 处理注册
		fmt.Println("处理注册逻辑")
		up := &process.UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessRegister(msg)
		if err != nil {
			fmt.Println("serverProcessLogin() 注册失败:", err)
			return
		}
	case message.SmsMsgType:
		fmt.Println("服务器接收到客户端发来群消息:", msg)
		smsProcess := process.SmsProcess{}
		smsProcess.SendGroupMsg(msg)
	case message.CtocMsgType:
		fmt.Println("服务器接收到客户端发来的私信:", msg)
		smsProcess := process.SmsProcess{}
		smsProcess.SendMsgToSpecifiedUser(msg)
	default:
		fmt.Println("不确定的消息类型.")
	}
	return
}

func (p *Processor) pkgReading() error{
	tf := &utils.Transfer{
		Conn: p.Conn,
	}
	for {
		//
		msg, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("readMsg(conn) 出错:", err)
			return err
		}
		fmt.Println("Msg:", msg)

		err = p.serverProcessMsg(&msg)
		if err != nil {
			return err
		}
	}
}