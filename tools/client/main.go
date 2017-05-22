package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/fwhappy/util"
)

var (
	host = flag.String("host", "0.0.0.0", "server host")
	port = flag.String("port", "38438", "server port")
	conn *net.TCPConn
	id   string
)

func init() {
	flag.Parse()
}

func main() {
	remote := *host + ":" + *port
	tcpAddr, err := net.ResolveTCPAddr("tcp", remote)
	if err != nil {
		fmt.Println("Error:ResolveTCPAddr:", err.Error())
		return
	}
	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Error:DialTCP:", err.Error())
		return
	}
	defer conn.Close()

	showClientDebug("connect at:%v", remote)

	showUsage()

	// 定义接收消息的协程
	go onRecived()

	// 控制台接收输入
	onInput()
}

func showUsage() {
	fmt.Println("------------------------------------------------")
	fmt.Println("usage:")
	fmt.Println("[握手]1.userId.nickname.extra1.extra2")
	fmt.Println("[下线]5")
	fmt.Println("[加入房间]4.1.roomId")
	fmt.Println("[退出房间]4.3.roomId")
	fmt.Println("[私人消息]4.101.userId.content")
	fmt.Println("[私人通知]4.103.userId.content")
	fmt.Println("[房间消息]4.105.roomId.content")
	fmt.Println("[房间通知]4.107.roomId.content")
	fmt.Println("[广播消息]4.109.content")
	fmt.Println("[广播通知]4.111.content")
	fmt.Println("------------------------------------------------")
}

// 显示客户端错误
func showClientError(a string, b ...interface{}) {
	fmt.Println("[", util.GetTimestamp(), "]", "[ERROR]", fmt.Sprintf(a, b...))
}

// 显示客户端调试信息
// 显示客户端错误
func showClientDebug(a string, b ...interface{}) {
	fmt.Println("[", util.GetTimestamp(), "]", "[DEBUG]", fmt.Sprintf(a, b...))
}
