package router

import (
	"goim/core"
	"net"

	"github.com/fwhappy/protocal"

	userController "goim/controller/user"
)

// Dispatch 转化客户端的请求
// 用户连接成功之后，userId，需要回传
func Dispatch(userId *string, conn *net.TCPConn, impacket *protocal.ImPacket, c chan int) {
	// 解析数据包
	packageId := impacket.GetPackage()
	switch packageId {
	case protocal.PACKAGE_TYPE_HANDSHAKE: // 握手
		// 记录当前连接的userId, 这样客户端无需每次传给服务端
		*userId = userController.HandShake(conn, impacket)
	case protocal.PACKAGE_TYPE_HANDSHAKE_ACK: // 握手成功
		userController.HandShakeAck(*userId, conn, impacket)
		c <- 1
	case protocal.PACKAGE_TYPE_HEARTBEAT: // 心跳
		userController.HeartBeat(*userId, conn, impacket)
	case protocal.PACKAGE_TYPE_DATA: // 数据包
		// 数据包路由分发
		routerData(*userId, conn, impacket)
	case protocal.PACKAGE_TYPE_KICK: // 下线
		// 直接调用defer触发的用户退出
		userController.Logout(*userId, conn, impacket)
	case protocal.PACKAGE_TYPE_SYSTEM: // 系统
		// game.SystemHandlerAction(conn, impacket)
	default:
		core.Logger.Error("未支持的数据包id:%d", packageId)
	}
}

// 数据包分发
func routerData(userId string, conn *net.TCPConn, impacket *protocal.ImPacket) {

}