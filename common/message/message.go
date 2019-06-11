package message

// 定义消息类型常量
const (
	LoginMsgType			= "LoginMsg"
	LoginResMsgType			= "LoginResMsg"
	RegisterMsgType			= "RegisterMsg"
	RegisterResMsgType		= "RegisterResMsg"
	NotifyUserStatusMsgType	= "NotifyUserStatusMsg"
	SmsMsgType				= "SmsMsg"
	CtocMsgType				= "CtocMsg"
)

// 定义用户状态常量
const (
	UserOffLine = iota
	UserOnline
	UserBusy
)

// 服务器与客户端之间互发消息的结构体
type Message struct {
	Type string  `json:"type"` // 消息类型
	Data string  `json:"data"` // 消息内容
}

// 用户登录消息体
type LoginMsg struct {
	UserId UserIdType `json:"userId"`		// 用户id
	UserPwd string `json:"userPwd"`	// 用户密码
	// UserName string `json:"userName"`	//用户名
}

// 服务器回复用户登录的消息
type LoginResMsg struct {
	Code int `json:"code"` // 返回状态码, 500 表示该用户还没有注册，200 表示登录成功
	OnlineUsers []User			// 保存在线用户id
	Error string `json:"error"` // 返回的错误信息
	User
}

type RegisterMsg struct {
	User
}

type RegisterResMsg struct {
	UserId UserIdType
	Code int `json:"code"`	// 返回状态码 400, 表示该用户已经存在，200 表示注册成功
	Error string `json:"error"`	// 返回错误信息
}

// 配合服务器推送用户状态变化的消息结构体
type NotifyUserStatusMsg struct {
	User
}

// 发送消息的结构体
type SmsMsg struct {
	User
	Context string `json:"context"`
}

// 点对点发送消息结构体
type CtocMsg struct {
	From User
	To User
	Context string `json:"context"`
}