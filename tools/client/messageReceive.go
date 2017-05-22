package main

import (
	"goim/config"
	"goim/core/msgpack"
	"io"
	"os"

	"strconv"

	"github.com/fwhappy/protocal"
	"github.com/fwhappy/util"
)

// 接受服务端消息
func onRecived() {
	for {
		// 读取包内容
		impacket, err := protocal.ReadPacket(conn)

		// 检查解析错误
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				showClientError("disconnected")
				os.Exit(0)
			} else {
				// 协议解析错误
				showClientError(err.Error())
			}
			break
		}

		body, _ := msgpack.Unmarshal(impacket.GetBody())
		switch impacket.GetPackageType() {
		case protocal.PACKAGE_TYPE_HANDSHAKE:
			s2cHandshake(body)
		case protocal.PACKAGE_TYPE_HANDSHAKE_ACK:
			s2cHandshakeAck()
		case protocal.PACKAGE_TYPE_HEARTBEAT:
			s2cHeartbeat()
		case protocal.PACKAGE_TYPE_DATA:
			// 收到握手ACK，开启心跳
			onRecivedData(impacket)
		}
	}
}

// 客户端回应收到的服务端数据
func onRecivedData(impacket *protocal.ImPacket) {
	body, err := msgpack.Unmarshal(impacket.GetBody())
	if err != nil {
		showClientError("解析服务器消息错误")
		return
	}

	s2cResult := body.JsonGetJsonMap("s2c_result")
	if code, _ := s2cResult.JsonGetInt("code"); code < 0 {
		errmsg, _ := s2cResult.JsonGetString("msg")
		showClientError("收到服务器返回的错误,messageId:%v, code:%v, errmsg:%v", impacket.GetMessageId(), code, errmsg)
		return
	}

	// 消息号
	messageNumber := impacket.GetMessageNumber()

	switch int(impacket.GetMessageId()) {
	case config.MESSAGE_ID_JOIN_ROOM_RESPONSE:
		s2cJoinRoomResponse(messageNumber, body)
	case config.MESSAGE_ID_JOIN_ROOM_PUSH:
		s2cJoinRoomPush(body)
	case config.MESSAGE_ID_QUIT_ROOM_RESPONSE:
		s2cQuitRoomResponse(messageNumber, body)
	case config.MESSAGE_ID_QUIT_ROOM_PUSH:
		s2cQuitRoomPush(body)
	case config.MESSAGE_ID_PRIVATE_MESSAGE_RESPONSE:
		s2cPrivateMessageResponse(messageNumber, body)
	case config.MESSAGE_ID_PRIVATE_MESSAGE_NOTIFY:
		s2cPrivateMessagePush(body)
	case config.MESSAGE_ID_ROOM_MESSAGE_RESPONSE:
		s2cRoomMessageResponse(messageNumber, body)
	case config.MESSAGE_ID_ROOM_MESSAGE_NOTIFY:
		s2cRoomMessagePush(body)
	case config.MESSAGE_ID_BROADCAST_MESSAGE_RESPONSE:
		s2cBroadcastMessageResponse(messageNumber, body)
	case config.MESSAGE_ID_BROADCAST_MESSAGE_NOTIFY:
		s2cBroadcastMessagePush(body)
	default:
		showClientError("为支持的服务端消息:%v", impacket.GetMessageId())
	}
}

// 加入房间的回应
func s2cJoinRoomResponse(messageNumber uint32, body util.JsonMap) {
	roomId, _ := body.JsonGetString("room_id")
	showClientDebug("加入房间成功, roomId:%v, number:%v", roomId, messageNumber)
}

// 其他人加入房间的推送
func s2cJoinRoomPush(body util.JsonMap) {
	roomId, _ := body.JsonGetString("room_id")
	showClientDebug("用户加入房间, roomId:%v, 用户信息:%v", roomId, body.JsonGetJsonMap("user"))
}

// 退出房间回应
func s2cQuitRoomResponse(messageNumber uint32, body util.JsonMap) {
	roomId, _ := body.JsonGetString("room_id")
	showClientDebug("退出房间成功, roomId:%v, number:", roomId, messageNumber)
}

// 退出房间通知
func s2cQuitRoomPush(body util.JsonMap) {
	roomId, _ := body.JsonGetString("room_id")
	showClientDebug("用户退出房间, roomId:%v, 用户信息:%v", roomId, body.JsonGetJsonMap("user"))
}

// 私聊消息回应
func s2cPrivateMessageResponse(messageNumber uint32, body util.JsonMap) {
	showClientDebug("私聊消息发送成功, number:%v", messageNumber)
}

// 推送私聊消息
func s2cPrivateMessagePush(body util.JsonMap) {
	content, _ := body.JsonGetString("content")
	showClientDebug("[private]%v, 用户信息:%v", content, body.JsonGetJsonMap("user"))
}

// 房间消息回应
func s2cRoomMessageResponse(messageNumber uint32, body util.JsonMap) {
	roomId, _ := body.JsonGetString("room_id")
	showClientDebug("房间消息发送成功, roomId:%v, number:%v", roomId, messageNumber)
}

// 收到房间消息
func s2cRoomMessagePush(body util.JsonMap) {
	roomId, _ := body.JsonGetString("room_id")
	content, _ := body.JsonGetString("content")
	showClientDebug("[room]%v, roomId:%v, 用户信息:%v", content, roomId, body.JsonGetJsonMap("user"))
}

// 广播消息回应
func s2cBroadcastMessageResponse(messageNumber uint32, body util.JsonMap) {
	showClientDebug("广播消息发送成功, number:%v", messageNumber)
}

// 推送广播消息
func s2cBroadcastMessagePush(body util.JsonMap) {
	content, _ := body.JsonGetString("content")
	showClientDebug("[broadcast]%v, 用户信息:%v", content, body.JsonGetJsonMap("user"))
}

// 收到服务端的握手消息
func s2cHandshake(body util.JsonMap) {
	interval, _ := body.JsonGetString("heartbeat_interval")
	heartbeatInterval, _ = strconv.Atoi(interval)
	go c2sHandShakeAck()
	showClientDebug("receive handshake")
}

// 收到服务端的握手回应消息
func s2cHandshakeAck() {
	// 收到握手ACK，开启心跳
	go c2sHeartBeat()
	showClientDebug("receive handshakeAck")
}

// 收到服务端的心跳
func s2cHeartbeat() {
	showClientDebug("receive heartbeah")
}
