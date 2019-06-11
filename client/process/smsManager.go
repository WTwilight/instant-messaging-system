package process

import (
	"encoding/json"
	"fmt"
	"github.com/WTwilight/instant-messaging-system/common/message"
)

// 显示群消息
func showGroupMsg(msg *message.Message) {

	var smsMsg message.SmsMsg
	err := json.Unmarshal([]byte(msg.Data), &smsMsg)
	if err != nil {
		fmt.Println("showGroupMsg() 中 json.Unmarshal([]byte(msg.Data), &smsMsg) 错误:", err)
	}
	// 显示信息
	fmt.Printf("来自ID: %d[%s] 的群发消息: %s\n\n", smsMsg.User.UserId, smsMsg.User.UserName, smsMsg.Context)
}

// 显示私信消息
func showCtocMsg(msg *message.Message)  {

	var ctocMsg message.CtocMsg
	err := json.Unmarshal([]byte(msg.Data), &ctocMsg)
	if err != nil {
		fmt.Println("showCtocMsg() 中 json.Unmarshal([]byte(msg.Data), &ctocMsg) 错误:", err)
	}

	fmt.Printf("来自ID: %s[%d] 的私信: %s\n\n", ctocMsg.From.UserName, ctocMsg.From.UserId, ctocMsg.Context)
}