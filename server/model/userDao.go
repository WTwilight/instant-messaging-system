package model

import (
	"encoding/json"
	"fmt"
	"github.com/WTwilight/instant-messaging-system/common/message"
	"github.com/gomodule/redigo/redis"
)

// Dao: Data Access Object

// 初始化一个全局 UserDao 实例，用于操作redis
var UserDaoHandle *UserDao

// 定义一个UserDao结构体
// 完成对User结构体的各种操作

type UserDao struct {
	pool *redis.Pool
}

// 使用工厂模式，创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	return &UserDao{
		pool: pool,
	}
}

// 1. 根据用户id, 返回一个User 实例和err
func (ud UserDao) getUserById(conn redis.Conn, id message.UserIdType) (user *message.User, err error) {
	userMsg, err := redis.String(conn.Do("HGET", "users", id))
	if err != nil {
		if err == redis.ErrNil {
			err = ErrorUserNotExist
		}
		return
	}

	user = &message.User{}
	err = json.Unmarshal([]byte(userMsg), user)
	if err != nil {
		fmt.Println("getUserById()中json.Unmarshal()错误:", err)
		return
	}
	return
}

// 2. 登录
func (ud *UserDao) Login(userId message.UserIdType, userPwd string) (user *message.User, err error) {
	conn := ud.pool.Get()
	defer func() {
		if conn.Close() != nil {
			fmt.Println("Login()关闭redis连接失败")
		}
	}()

	user, err = ud.getUserById(conn, userId)
	if err != nil {
		return
	}

	if user.UserPwd != userPwd {
		err = ErrorUserWrongPwd
		return
	}
	return
}

// 2. 注册
func (ud *UserDao) Register(user *message.User) (err error) {
	conn := ud.pool.Get()
	defer func() {
		if conn.Close() != nil {
			fmt.Println("Register()关闭redis连接失败")
		}
	}()

	_, err = ud.getUserById(conn, user.UserId)
	if err == nil {
		// 如果 没有出错，说明从数据库中成功取出了userId, 说明用户已经存在，不能再使用此id注册
		err = ErrorUserExist
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Register() 中 json.Marshal(user) 序列化失败:", err)
		return
	}

	_, err = conn.Do("HSET", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存用户注册信息错误:", err)
		return
	}

	_, err = conn.Do("LPUSH", "useridlist", user.UserId)
	if err != nil {
		fmt.Println("保存用户id list 错误:", err)
		return
	}
	return
}

func (ud *UserDao) GenerateNewUserId() (id message.UserIdType, err error) {
	conn := ud.pool.Get()
	defer func() {
		if conn.Close() != nil {
			fmt.Println("Register()关闭redis连接失败")
		}
	}()

	// 获取 redis 中 useridlist 这个list的长度，如果useridlist不存在或者为空，都会返回0
	offset, err := redis.Int64(conn.Do("LLEN", "useridlist"))
	if err != nil {
		fmt.Println("GetLargestIdFromUserIdList() 中 redis.Int64(conn.Do()) 错误:", err)
		return
	}

	id = message.UserIdStart + message.UserIdType(offset)
	return
}