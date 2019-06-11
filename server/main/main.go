package main

import (
	"fmt"
	"net"
)

func main() {

	fmt.Println("服务器[新结构]在 8889 端口监听...")
	listen, err := net.Listen("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Listen() 错误:", err)
		return
	}

	for  {
		fmt.Println("等待客户端连接...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept() 出错:", err)
		}

		// 连接成功则启动一个协程和客户保持通讯...
		go processConn(conn)
	}
}

// 处理和客户端的通讯
func processConn(conn net.Conn)  {
	// 延时关闭连接
	defer func() {
		if conn.Close() != nil {
			fmt.Println("服务器端关闭客户端连接失败.")
		}
	}()
	//
	processor := &Processor{
		Conn: conn,
	}
	err := processor.pkgReading()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误:", err)
		// time.Sleep(time.Second * 3)
		return
	}
}