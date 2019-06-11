package message

type UserIdType uint64

const (
	UserIdStart UserIdType = 100000
	UserIdEnd   UserIdType = 999999
)

type User struct {
	UserStatus	int		`json:"userStatus"`	// 用户状态
	UserAge		int		`json:"userAge"`
	UserId UserIdType	`json:"userId"`
	UserPwd string		`json:"userPwd"`
	UserName string		`json:"userName"`
	UserGender string	`json:"userGender"`
}