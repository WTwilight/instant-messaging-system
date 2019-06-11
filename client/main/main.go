package main

import (
	"fmt"
	"github.com/WTwilight/instant-messaging-system/client/process"
	"github.com/WTwilight/instant-messaging-system/common/message"
)

// 定义两个变量，一个保存用户id,用保存用户密码
var (
	userId message.UserIdType
	userPwd string
	userName string
	user message.User
)

func showMainMenu() {
	fmt.Println("----------欢迎登录多人聊天室----------")
	fmt.Println("            1. 登录聊天室")
	fmt.Println("            2. 注册用户")
	fmt.Println("            3. 退出系统")
	fmt.Println("请选择(1-3):")
}

func receiveUserLoginInput(userId *message.UserIdType, userPwd *string)  {
	fmt.Println("输入用户id:")
	fmt.Scanf("%d\n", userId)

	fmt.Println("输入密码:")
	fmt.Scanf("%s\n", userPwd)
}

func receiveUserRegisterInput(user *message.User) {
	fmt.Println("输入用户名:")
	fmt.Scanf("%s\n", &user.UserName)

	fmt.Println("输入密码:")
	fmt.Scanf("%s\n", &user.UserPwd)

	fmt.Println("输入性别:")
	fmt.Scanf("%s\n", &user.UserGender)

	fmt.Println("输入年龄:")
	fmt.Scanf("%d\n", &user.UserAge)
}

func main() {
	//接收用户的选择
	var key int

	// 判断是否还继续显示菜单
	var loop = true

	for  {
		if !loop {
			break
		}

		showMainMenu()
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			// 登录聊天室
			receiveUserLoginInput(&userId, &userPwd)
			up := process.UserProcess{}
			_ = up.Login(userId, userPwd)

		case 2:
			// 注册聊天室
			// 1. 获取用户输入
			receiveUserRegisterInput(&user)
			fmt.Println("用户输入的注册信息:", user)
			// 2. 调用UserProcess， 完成注册请求
			up := process.UserProcess{}
			_ = up.Register(&user)
			//loop = false
		case 3:
			fmt.Printf("退出系统\n")
			loop = false
		default:
			fmt.Printf("输入有误，重新输入.\n")
		}
	}
}
