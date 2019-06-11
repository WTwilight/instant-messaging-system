package process

import (
	"encoding/json"
	"fmt"
	"github.com/WTwilight/instant-messaging-system/client/utils"
	"github.com/WTwilight/instant-messaging-system/common/message"
)

type SmsProcess struct {
	//
}

func (sp *SmsProcess) SendGroupMsg(content string) (err error) {

	// 1. 创建一个 SmsMsg 实例
	smsMsg := message.SmsMsg{
		Context: content,
	}
	smsMsg.User = CurrentUser.User

	// 2. 序列化 smsMsg
	data,err := json.Marshal(smsMsg)
	if err != nil {
		fmt.Println("SendGroupMsg() 中 json.Marshal(smsMsg) 错误:", err)
		return
	}

	// 3. 创建一个 message.Message 消息实例
	msg := message.Message{
		Data: string(data),
		Type: message.SmsMsgType,
	}

	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("SendGroupMsg() 中 json.Marshal(msg) 错误:", err)
		return
	}

	// 4. 将msg发送给服务器
	tf := &utils.Transfer{
		Conn: CurrentUser.Conn,
	}

	// 5. 发送 msg
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMsg() 中 tf.WritePkg(data) 错误:", err.Error())
	}
	return
}

func (sp *SmsProcess) SendMsgByUserID(userId message.UserIdType) (err error) {
	// 给userId发送消息
	// 1. 从onlineUsers中找到对应的用户
	user, err := getUserInfoById(userId)
	if err != nil {
		fmt.Println(err)
		return
	}

	var context string
	fmt.Printf("想对 %s[%d] 说(输入要说的话):", user.UserName, user.UserId)
	_, _ = fmt.Scanf("%s\n", &context)

	ctocMsg := &message.CtocMsg{
		From: CurrentUser.User,
		To: *user,
		Context: context,
	}

	data, err := json.Marshal(ctocMsg)
	if err != nil {
		fmt.Println("SendMsgByUserID() 中 json.Marshal(ctocMsg) 错误:", err)
		return
	}

	msg := &message.Message{
		Type: message.CtocMsgType,
		Data: string(data),
	}

	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("SendMsgByUserID() 中 json.Marshal(msg) 错误:", err)
		return
	}

	tf := utils.Transfer{
		Conn: CurrentUser.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Printf("给 %s[%d] 发送消息失败: %v\n", user.UserName, user.UserId, err)
	}

	return
}