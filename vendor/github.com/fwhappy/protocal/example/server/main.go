package main

import (
	"fmt"
	"io"
	"net"
	"protocal"
	"runtime"
	"time"
)

func main() {
	remote := "0.0.0.0:38438"
	tcpAddr, resolveErr := net.ResolveTCPAddr("tcp", remote)
	if resolveErr != nil {
		fmt.Println("listenAndServe ResolveTCPAddr Error:", resolveErr.Error())
		return
	}
	tcpListener, listenErr := net.ListenTCP("tcp", tcpAddr)
	if listenErr != nil {
		fmt.Println("listenAndServe ListenTCP Error:", listenErr.Error())
		return
	}
	fmt.Println("[", time.Now().Format("2006-01-02 15:04:05"), "]", "success listen:", remote)

	// 监听连接事件
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println("tcpListener.AcceptTCP:", err)
			continue
		}

		// 客户端连接成功，开启新的协程，监听客户端消息
		go serve(tcpConn)
	}
}

// serve 监听客户端连接
func serve(conn *net.TCPConn) {
	defer func() {
		// 捕获异常
		if err := recover(); err != nil {
			stack := make([]byte, 1024)
			stack = stack[:runtime.Stack(stack, true)]

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Println("[", timestamp, "]", "catchPanic:", err)
			fmt.Println("[", timestamp, "]", "stack:", string(stack))
		}
	}()

	for {
		// 读取包内容
		impacket, err := protocal.ReadPacket(conn)
		// 检查解析错误
		if err != nil {
			switch err {
			case io.EOF:
				// 关闭退出
				fmt.Println("User disconnected, remote:", conn.RemoteAddr().String())
			case io.ErrUnexpectedEOF:
				fmt.Println("unexpected EOF, remote:", conn.RemoteAddr().String())
			default:
				// 协议解析错误
				fmt.Println(err.Error())
			}
			break
		}

		fmt.Println("收到客户端消息")
		fmt.Println("包类型:", impacket.GetPackageType())
		fmt.Println("包体长度:", impacket.GetMessageLength())
		fmt.Println("包体内容:", string(impacket.GetBody()))

		time.Sleep(time.Second)

		// 恢复一条消息
		message := protocal.NewImMessage(uint16(1), protocal.MSG_TYPE_PUSH, uint32(1), []byte("message from server."))
		packet := protocal.NewImPacket(protocal.PACKAGE_TYPE_DATA, message)
		packet.Send(conn)
	}
}
