package user

import (
	"goim/response"
	"net"

	"github.com/fwhappy/protocal"

	userService "goim/service/user"
)

// HandShake 用户握手
func HandShake(conn *net.TCPConn, impacket *protocal.ImPacket) string {
	id, err := userService.HandShake(conn, impacket)
	if err != nil {
		// 握手失败
		response.GetError(protocal.PACKAGE_TYPE_HANDSHAKE, err, nil).Send(conn)
	}
	return id
}

// HandShakeAck 握手回应
func HandShakeAck(id string, conn *net.TCPConn, impacket *protocal.ImPacket) {
	if err := userService.HandShakeAck(id, impacket); err != nil {
		response.GetError(protocal.PACKAGE_TYPE_HANDSHAKE_ACK, err, nil).Send(conn)
	}
}

// HeartBeat 用户心跳
func HeartBeat(id string, conn *net.TCPConn, impacket *protocal.ImPacket) {
	userService.HeartBeat(id)
}

// Logout 用户请求退出
func Logout(id string, conn *net.TCPConn, impacket *protocal.ImPacket) {
	if err := userService.Logout(id); err != nil {
		response.GetError(protocal.PACKAGE_TYPE_KICK, err, nil).Send(conn)
	}
}
