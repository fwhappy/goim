package main

import (
	"flag"
	"fmt"
	"goim/config"
	"net"

	"github.com/fwhappy/protocal"
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
	fmt.Println("[握手]", protocal.PACKAGE_TYPE_HANDSHAKE, ".userId.nickname.extra1.extra2")
	fmt.Println("[下线]", protocal.PACKAGE_TYPE_KICK)
	fmt.Println("[加入房间]", protocal.PACKAGE_TYPE_DATA, ".", config.MESSAGE_ID_JOIN_ROOM_REQUEST, ".roomId")
	fmt.Println("[退出房间]", protocal.PACKAGE_TYPE_DATA, ".", config.MESSAGE_ID_QUIT_ROOM_REQUEST, ".roomId")
	fmt.Println("[私人消息]", protocal.PACKAGE_TYPE_DATA, ".", config.MESSAGE_ID_PRIVATE_MESSAGE_REQUEST, ".userId.content")
	// fmt.Println("[私人通知]", protocal.PACKAGE_TYPE_DATA, ".", config.MESSAGE_ID_PRIVATE_MESSAGE_NOTIFY, ".userId.content")
	fmt.Println("[房间消息]", protocal.PACKAGE_TYPE_DATA, ".", config.MESSAGE_ID_ROOM_MESSAGE_REQUEST, ".roomId.content")
	// fmt.Println("[房间通知]", protocal.PACKAGE_TYPE_DATA, ".", config.MESSAGE_ID_ROOM_MESSAGE_NOTIFY, ".roomId.content")
	fmt.Println("[广播消息]", protocal.PACKAGE_TYPE_DATA, ".", config.MESSAGE_ID_BROADCAST_MESSAGE_REQUEST, ".content")
	// fmt.Println("[广播通知]", protocal.PACKAGE_TYPE_DATA, ".", config.MESSAGE_ID_BROADCAST_MESSAGE_NOTIFY, ".content")
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
