package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/WTwilight/instant-messaging-system/common/message"
	"io"
	"net"
)

// 将这些方法关联到结构体中
type Transfer struct {
	Conn net.Conn
	Buf [8192]byte
}

func (tf *Transfer)ReadPkg() (msg message.Message, err error) {

	//buf := make([]byte, 8192)
	fmt.Println("读取客户端发送的数据...")
	// conn.Read 只有在 conn 没有被关闭的情况下，才会阻塞
	// 如果客户端关闭了 conn 则，就不会阻塞.
	_, err = tf.Conn.Read(tf.Buf[:4])
	if err != nil {
		if err == io.EOF {
			fmt.Println("服务端关闭了连接.退出服务...")
			return
		} else {
			fmt.Println("conn.Read() 出错:", err)
		}
	}
	fmt.Println("读到bytesBuf:", tf.Buf[:4])

	// 根据 buf[:4] 转成一个 uint32 类型
	pkgLen := binary.BigEndian.Uint32(tf.Buf[0:4])

	fmt.Println("消息的长度:", pkgLen)
	// 根据 msgLen 读取消息内容
	n, err := tf.Conn.Read(tf.Buf[0:pkgLen])
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("conn.Read 出错:", err)
		return
	}

	// 把 buf 中的 msg 信息反序列化成 Message 结构体
	if json.Unmarshal(tf.Buf[:pkgLen], &msg) != nil {
		fmt.Println("反序列失败:", err)
		return
	}
	return
}

func (tf *Transfer)WritePkg(data []byte) (err error) {

	// 先发送data的长度 给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))

	//buf := make([]byte, 4)
	binary.BigEndian.PutUint32(tf.Buf[:4], pkgLen)
	// 发送长度
	n, err := tf.Conn.Write(tf.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(lenBytes) 失败:", err)
		return
	}

	//发送 data 本身
	n, err = tf.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(data) 失败:", err)
		return
	}
	return
}