package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"protocal"
	"strings"
)

func main() {
	remote := "0.0.0.0:38438"
	tcpAddr, err := net.ResolveTCPAddr("tcp", remote)
	if err != nil {
		fmt.Println("Error:ResolveTCPAddr:", err.Error())
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Error:DialTCP:", err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("connect at: ", remote)

	// 定义接收消息的协程
	go onMessageRecived(conn)

	// 控制台接收输入
	onMessageInput(conn)
}

// 接受控制台输入，并发送消息给服务器
func onMessageInput(conn *net.TCPConn) {
	for {
		var msg string
		msgReader := bufio.NewReader(os.Stdin)
		msg, _ = msgReader.ReadString('\n')
		msg = strings.TrimSuffix(msg, "\n")

		// 跳过空数据包
		if len(msg) == 0 {
			continue
		} else if msg == "quit" || msg == "exit" {
			break
		} else if msg == "cancel" {
			conn.Close()
		}

		packageType := protocal.PACKAGE_TYPE_DATA
		messageId := uint16(1)
		messageType := protocal.MSG_TYPE_REQUEST
		messageNumber := uint32(10)

		message := protocal.NewImMessage(messageId, messageType, messageNumber, []byte(msg))
		packet := protocal.NewImPacket(packageType, message)
		packet.Send(conn)
	}
}

func onMessageRecived(conn *net.TCPConn) {
	for {
		// 读取包内容
		impacket, err := protocal.ReadPacket(conn)

		// 检查解析错误
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				fmt.Println("disconnected")
				os.Exit(0)
			} else {
				// 协议解析错误
				fmt.Println(err.Error())
			}
			break
		}

		fmt.Println("收到服务端消息")
		fmt.Println("包类型:", impacket.GetPackageType())
		fmt.Println("消息长度:", impacket.GetMessageLength())
		fmt.Println("消息id:", impacket.GetMessageId())
		fmt.Println("消息类型:", impacket.GetMessageType())
		fmt.Println("消息编号:", impacket.GetMessageNumber())
		fmt.Println("消息正文:", string(impacket.GetBody()))
	}
}
