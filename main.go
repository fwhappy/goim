package main

import (
	"flag"
	"goim/core"
	"goim/hall"
	"goim/router"
	hallService "goim/service/hall"
	userService "goim/service/user"
	"io"
	"net"
	"runtime"

	"github.com/fwhappy/protocal"
	"github.com/fwhappy/util"
)

var (
	// 监听端口
	port = flag.String("p", "38438", "listen port")
	// 配置文件夹, 最好是绝对路劲, /Home/xh/goim/etc/local
	confDir = flag.String("confDir", "etc/local", "config dir path")
)

func init() {
	// 解析url参数
	flag.Parse()
}

func main() {
	defer util.RecoverPanic()

	// // 初始化基础配置
	core.LoadAppConfig(core.GetConfigFile("app.toml", *confDir))
	// // 初始化日志配置
	core.LoadLoggerConfig(core.GetConfigFile("log.toml", *confDir))
	defer core.Logger.Flush()

	// 开启监听端口 & 提供服务
	listenAndServe()
}

// listenAndServe 开始监听并提供服务
func listenAndServe() {
	listenRemote := "0.0.0.0:" + *port
	tcpAddr, resolveErr := net.ResolveTCPAddr("tcp", listenRemote)
	if resolveErr != nil {
		core.Logger.Error("listenAndServe ResolveTCPAddr Error:%s", resolveErr.Error())
		return
	}
	tcpListener, listenErr := net.ListenTCP("tcp", tcpAddr)
	if listenErr != nil {
		core.Logger.Error("listenAndServe ListenTCP Error:%s", listenErr.Error())
		return
	}
	core.Logger.Info("server lisen at: " + listenRemote)

	// 监听连接事件
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			core.Logger.Error("tcpListener.AcceptTCP: %s.", err)
			continue
		}

		// 客户端连接成功，开启新的协程，监听客户端消息
		go serve(tcpConn)
	}
}

func serve(conn *net.TCPConn) {
	core.Logger.Debug("New User connected: %s.", conn.RemoteAddr().String())

	// 记录当前连接的用户id
	var userId string

	// 当前用户连接成功或出错退出的消息
	c := make(chan int, 2)
	// 定义异常捕捉
	defer func() {
		// 捕获异常
		if err := recover(); err != nil {
			core.Logger.Error("serve defer recover error: %s.", err)

			stack := make([]byte, 1024)
			stack = stack[:runtime.Stack(stack, true)]
			core.Logger.Debug("stack:\n%s", string(stack))
			core.Logger.Debug("defer disconnected: %s.", conn.RemoteAddr().String())
		}

		// 断开用户连接
		if userId != "" {
			if u, online := hall.UserSet.Get(userId); online && u.Conn == conn {
				userService.KickUser(u)
			}
		} else {
			conn.Close()
			c <- -1
		}
	}()

	// 检测用户连接之后，如果在规定时间内handshake成功，需要断开连接，防止无效的连接
	go hallService.ListenHandShakeTimeout(conn, c)

	// 解析消息
	for {
		// 读取包内容
		impacket, err := protocal.ReadPacket(conn)

		// 检查解析错误
		if err != nil {
			switch err {
			case io.EOF:
				// 关闭退出
				core.Logger.Debug("User disconnected, remote: %s.", conn.RemoteAddr().String())
			case io.ErrUnexpectedEOF:
				core.Logger.Debug("unexpected EOF, remote: %s.", conn.RemoteAddr().String())
			default:
				// 协议解析错误
				core.Logger.Error(err.Error())
			}
			break
		}

		router.Dispatch(&userId, conn, impacket, c)
	}
}
