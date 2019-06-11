package model

import "errors"

// 根据业务逻辑需要自定义错误
var (
	ErrorUserNotExist = errors.New("用户不存在")
	ErrorUserExist = errors.New("用户已存在")
	ErrorUserWrongPwd = errors.New("密码不正确")
)

