package main

import (
	"bufio"
	"goim/config"
	"goim/core/msgpack"
	"os"
	"protocal"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fwhappy/util"
)

// 消息号码生成器
type numberGenerator struct {
	value int
	mux   *sync.Mutex
}

var (
	mg                *numberGenerator
	heartbeatInterval int
)

// 生成一个消息号
func (mg *numberGenerator) getNumber() uint32 {
	mg.mux.Lock()
	defer mg.mux.Unlock()
	mg.value++
	return uint32(mg.value)
}

// 发送消息（非data类型）
func send(packageType uint8, body util.JsonMap) {
	var message []byte
	if body != nil {
		message = msgpack.Marshal(body)
	}
	protocal.NewImPacket(packageType, message).Send(conn)
}

// 发送data类型的消息
func sendData(messageId, messageType uint16, body util.JsonMap) {
	var message []byte
	if body != nil {
		message = protocal.NewImMessage(messageId, messageType, mg.getNumber(), msgpack.Marshal(body))
	}
	protocal.NewImPacket(protocal.PACKAGE_TYPE_DATA, message).Send(conn)
}

// 接受客户端输入
func onInput() {
	for {
		var input string
		input, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		if len(input) == 0 {
			continue
		} else if input == "quit" || input == "exit" {
			break
		} else if input == "close" {
			conn.Close()
		} else if input == "help" || input == "usage" {
			showUsage()
			continue
		}

		// 包类型
		var packageType uint8
		params := strings.Split(input, ".")
		p1, _ := strconv.Atoi(params[0])
		packageType = uint8(p1)
		switch packageType {
		case protocal.PACKAGE_TYPE_HANDSHAKE: // 握手
			c2sHandShake(params[1:]...)
		// case protocal.PACKAGE_TYPE_HANDSHAKE_ACK: // 握手完成
		// case protocal.PACKAGE_TYPE_HEARTBEAT: // 心跳
		case protocal.PACKAGE_TYPE_KICK: // 退出
			c2sLogout()
		case protocal.PACKAGE_TYPE_DATA: // 数据包
			onInputData(params[1:])
		default:
			showClientError("为支持的包类型:%v", packageType)
		}
	}
}

// 处理用户输入的input数据
func onInputData(args []string) {
	messageId := getParamsInt(0, args)
	switch messageId {
	case config.MESSAGE_ID_JOIN_ROOM_REQUEST:
		c2sJoinRoomRequest(args)
	case config.MESSAGE_ID_QUIT_ROOM_REQUEST:
		c2sQuitRoomRequset(args)
	case config.MESSAGE_ID_PRIVATE_MESSAGE_REQUEST:
		c2sPrivateMessageRequest(args)
	case config.MESSAGE_ID_PRIVATE_MESSAGE_NOTIFY:
		c2sPrivateMessageNotify(args)
	case config.MESSAGE_ID_ROOM_MESSAGE_REQUEST:
		c2sRoomMessageRequest(args)
	case config.MESSAGE_ID_ROOM_MESSAGE_NOTIFY:
		c2sRoomMessageNotify(args)
	case config.MESSAGE_ID_BROADCAST_MESSAGE_REQUEST:
		c2sBroadcastMessageRequest(args)
	case config.MESSAGE_ID_BROADCAST_MESSAGE_NOTIFY:
		c2sBroadcastMessageNotify(args)
	default:
		showClientError("为支持的消息id:%v", messageId)
	}
}

// 客户端加入房间请求
func c2sJoinRoomRequest(args []string) {
	roomId := getParams(1, args)
	if roomId == "" {
		showClientError("未输入房间号")
		return
	}
	body := util.JsonMap{"roomId": roomId}
	sendData(config.MESSAGE_ID_JOIN_ROOM_REQUEST, protocal.MSG_TYPE_REQUEST, body)
}

// 客户端退出房间请求
func c2sQuitRoomRequset(args []string) {
	roomId := getParams(1, args)
	if roomId == "" {
		showClientError("未输入房间号")
		return
	}
	body := util.JsonMap{"roomId": roomId}
	sendData(config.MESSAGE_ID_QUIT_ROOM_REQUEST, protocal.MSG_TYPE_REQUEST, body)
}

// 客户端私人消息请求
func c2sPrivateMessageRequest(args []string) {
	toId := getParams(1, args)
	content := getParams(2, args)
	if toId == "" || content == "" {
		showClientError("私聊对象和私聊内容不能为空,toId:%v, content:%v", toId, content)
	}
	body := util.JsonMap{"id": toId, "content": content}
	sendData(config.MESSAGE_ID_PRIVATE_MESSAGE_REQUEST, protocal.MSG_TYPE_REQUEST, body)
}

// 客户端私人消息通知
func c2sPrivateMessageNotify(args []string) {
	toId := getParams(1, args)
	content := getParams(2, args)
	if toId == "" || content == "" {
		showClientError("私聊对象和私聊内容不能为空,toId:%v, content:%v", toId, content)
	}
	body := util.JsonMap{"id": toId, "content": content}
	sendData(config.MESSAGE_ID_PRIVATE_MESSAGE_NOTIFY, protocal.MSG_TYPE_NOTIFY, body)
}

// 客户端房间消息请求
func c2sRoomMessageRequest(args []string) {
	toId := getParams(1, args)
	content := getParams(2, args)
	if toId == "" || content == "" {
		showClientError("房间id和内容不能为空,toId:%v, content:%v", toId, content)
	}
	body := util.JsonMap{"id": toId, "content": content}
	sendData(config.MESSAGE_ID_ROOM_MESSAGE_REQUEST, protocal.MSG_TYPE_REQUEST, body)
}

// 客户端房间消息通知
func c2sRoomMessageNotify(args []string) {
	toId := getParams(1, args)
	content := getParams(2, args)
	if toId == "" || content == "" {
		showClientError("房间id和内容不能为空,toId:%v, content:%v", toId, content)
	}
	body := util.JsonMap{"id": toId, "content": content}
	sendData(config.MESSAGE_ID_ROOM_MESSAGE_NOTIFY, protocal.MSG_TYPE_NOTIFY, body)
}

// 客户端广播房间消息请求
func c2sBroadcastMessageRequest(args []string) {
	content := getParams(1, args)
	if content == "" {
		showClientError("内容不能为空, content:%v", content)
	}
	body := util.JsonMap{"content": content}
	sendData(config.MESSAGE_ID_BROADCAST_MESSAGE_REQUEST, protocal.MSG_TYPE_REQUEST, body)
}

// 客户端广播消息通知
func c2sBroadcastMessageNotify(args []string) {
	content := getParams(1, args)
	if content == "" {
		showClientError("内容不能为空, content:%v", content)
	}
	body := util.JsonMap{"content": content}
	sendData(config.MESSAGE_ID_BROADCAST_MESSAGE_NOTIFY, protocal.MSG_TYPE_NOTIFY, body)
}

// 客户端向服务端发送握手协议
func c2sHandShake(args ...string) {
	userInfo := make(util.JsonMap)
	extra := make(util.JsonMap)
	userInfo["id"] = args[0]
	userInfo["nickname"] = args[1]
	extra["extra1"] = args[2]
	extra["extra2"] = args[3]
	userInfo["extra"] = extra
	id = args[0]
	send(protocal.PACKAGE_TYPE_HANDSHAKE, userInfo)

	showClientDebug("send handShake")
}

// client给server发送握手成功
func c2sHandShakeAck() {
	// 发送消息给服务器
	send(protocal.PACKAGE_TYPE_HANDSHAKE_ACK, nil)
	showClientDebug("send handShakeAck")
}

// client每3秒给server发送一个心跳消息
// 服务端如果超过6秒没有收到包，则认为客户端已离线
func c2sHeartBeat() {
	for {
		time.Sleep(time.Duration(heartbeatInterval) * time.Second)
		// 发送消息给服务器
		send(protocal.PACKAGE_TYPE_HEARTBEAT, nil)
		showClientDebug("send c2sHeartBeat")
	}
}

// 用户退出
func c2sLogout() {
	// 发送消息给服务器
	send(protocal.PACKAGE_TYPE_KICK, nil)
	showClientDebug("send c2sLogout")
}

func getParams(position int, args []string) string {
	if position > len(args) {
		return ""
	}
	return args[position]
}

func getParamsInt(position int, args []string) int {
	param := getParams(position, args)
	if param == "" {
		return 0
	}
	value, _ := strconv.Atoi(param)
	return value
}
