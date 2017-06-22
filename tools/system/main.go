package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/fwhappy/protocal"
	"github.com/fwhappy/util"
)

var (
	host   = flag.String("host", "0.0.0.0", "server host")
	port   = flag.String("port", "38438", "server port")
	secret = flag.String("secret", "3848438", "secret key")
	conn   *net.TCPConn
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
	fmt.Println("[用户信息]", protocal.PACKAGE_TYPE_SYSTEM, ".userId")
	fmt.Println("[房间信息]", protocal.PACKAGE_TYPE_SYSTEM, ".roomId")
	fmt.Println("[踢人]", protocal.PACKAGE_TYPE_SYSTEM, ".userId")
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
