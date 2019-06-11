package process

import (
	"encoding/json"
	"fmt"
	"github.com/WTwilight/instant-messaging-system/common/message"
	"github.com/WTwilight/instant-messaging-system/server/utils"
	"net"
)

type SmsProcess struct {
	//
}

func (sp *SmsProcess) SendMsgToOnlineUser(data []byte, conn net.Conn) (err error) {
	//
	tf := utils.Transfer{
		Conn: conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendMsgToOnlineUser() 中 tf.WritePkg(data) 错误:", err.Error())
	}
	fmt.Println("SendMsgToOnlineUser()服务器转发短消息成功")
	return
}

func (sp *SmsProcess) SendGroupMsg(msg *message.Message) {

	// 遍历 onlineUsers ，将消息转发出去

	var smsMsg message.SmsMsg
	err := json.Unmarshal([]byte(msg.Data), &smsMsg)
	if err != nil {
		fmt.Println("SendGroupMsg(msg *message.Message) 中 json.Unmarshal([]byte(msg.Data), &smsMsg) 错误:", err)
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("SendGroupMsg(msg *message.Message) 中 json.Marshal(msg) 错误:", err)
		return
	}

	for id, onlineUser := range userMgr.OnlineUsers {
		//
		if message.UserIdType(id) == smsMsg.User.UserId {
			continue
		}
		err = sp.SendMsgToOnlineUser(data, onlineUser.Conn)
		if err != nil {
			fmt.Println("SendGroupMsg(msg *message.Message) 中 SendMsgToOnlineUser(smsMsg.Context, onlineUser.Conn) 错误:", err)
		}
	}
}

func (sp *SmsProcess) SendMsgToSpecifiedUser(msg *message.Message) {
	//
	var ctocMsg message.CtocMsg
	if json.Unmarshal([]byte(msg.Data), &ctocMsg) != nil {
		fmt.Println("SendMsgToSpecifiedUser() 中 json.Unmarshal() 错误.")
		return
	}

	up, ok := userMgr.OnlineUsers[ctocMsg.To.UserId]
	if !ok {
		fmt.Println("用户不在线")
		// 回复发送者用户不在线，发送失败。
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("SendMsgToSpecifiedUser() 中 json.Marshal(msg) 错误:", err)
		return
	}

	tf := utils.Transfer{
		Conn: up.Conn,
	}
	if tf.WritePkg(data) != nil {
		fmt.Println("SendMsgToSpecifiedUser() 中 tf.WritePkg(data) 错误")
	}
}