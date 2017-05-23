package message

import (
	"goim/config"
	"goim/core/msgpack"
	"goim/hall"
	"goim/ierror"
	"goim/response"

	"github.com/fwhappy/protocal"
	"github.com/fwhappy/util"
)

// PrivateMessage 私聊消息
func PrivateMessage(id string, impacket *protocal.ImPacket) *ierror.Error {
	// 解析协议
	body, err := msgpack.Unmarshal(impacket.GetBody())
	if err != nil {
		return ierror.NewError(-100, "private message", err.Error())
	}

	// 消息内容
	content, _ := body.JsonGetString("content")
	if content != "" {
		if !checkAndFilterContent(content) {
			return ierror.NewError(-403)
		}
	}
	if content == "" {
		return ierror.NewError(-402)
	}

	// 私聊对象
	target, _ := body.JsonGetString("id")
	if target == "" {
		return ierror.NewError(-400)
	}

	user, online := hall.UserSet.Get(id)
	if !online {
		return ierror.NewError(-202)
	}
	u, online := hall.UserSet.Get(target)
	if !online {
		return ierror.NewError(-401, id)
	}

	// 回应一个成功
	userResponse := response.GetSuccessData(config.MESSAGE_ID_PRIVATE_MESSAGE_RESPONSE, protocal.MSG_TYPE_RESPONSE, impacket.GetMessageNumber(), nil)
	user.AppendMessage(userResponse)

	// 给对象发消息
	pushBody := util.JsonMap{"content": content, "user": user.Info}
	targetPush := response.GetSuccessData(config.MESSAGE_ID_PRIVATE_MESSAGE_PUSH, protocal.MSG_TYPE_PUSH, impacket.GetMessageNumber(), pushBody)
	u.AppendMessage(targetPush)

	return nil
}

// RoomMessage 房间消息
func RoomMessage(id string, impacket *protocal.ImPacket) *ierror.Error {
	// 解析协议
	body, err := msgpack.Unmarshal(impacket.GetBody())
	if err != nil {
		return ierror.NewError(-100, "private message", err.Error())
	}

	// 消息内容
	content, _ := body.JsonGetString("content")
	if content != "" {
		if !checkAndFilterContent(content) {
			return ierror.NewError(-403)
		}
	}
	if content == "" {
		return ierror.NewError(-402)
	}

	user, online := hall.UserSet.Get(id)
	if !online {
		return ierror.NewError(-202)
	}

	// 私聊对象
	target, _ := body.JsonGetString("room_id")
	if target == "" {
		return ierror.NewError(-400)
	}
	r, exists := hall.RoomSet.Get(target)
	if !exists {
		return ierror.NewError(-404)
	}

	// 判断用户是否在房间内
	if !r.IsUserInRoom(id) {
		return ierror.NewError(-303, id, target)
	}

	// 回应一个成功
	userResponse := response.GetSuccessData(config.MESSAGE_ID_ROOM_MESSAGE_RESPONSE, protocal.MSG_TYPE_RESPONSE, impacket.GetMessageNumber(), nil)
	user.AppendMessage(userResponse)

	// 给房间成员发消息
	pushBody := util.JsonMap{"content": content, "user": user.Info}
	targetPush := response.GetSuccessData(config.MESSAGE_ID_ROOM_MESSAGE_PUSH, protocal.MSG_TYPE_PUSH, impacket.GetMessageNumber(), pushBody)
	hall.SendRoomMessage(r, targetPush)

	return nil
}

// BroadcastMessage 广播消息
func BroadcastMessage(id string, impacket *protocal.ImPacket) *ierror.Error {
	// TODO 检测用户的发送频率

	// 解析协议
	body, err := msgpack.Unmarshal(impacket.GetBody())
	if err != nil {
		return ierror.NewError(-100, "private message", err.Error())
	}

	// 消息内容
	content, _ := body.JsonGetString("content")
	if content != "" {
		if !checkAndFilterContent(content) {
			return ierror.NewError(-403)
		}
	}
	if content == "" {
		return ierror.NewError(-402)
	}

	user, online := hall.UserSet.Get(id)
	if !online {
		return ierror.NewError(-202)
	}

	// 回应一个成功
	userResponse := response.GetSuccessData(config.MESSAGE_ID_BROADCAST_MESSAGE_RESPONSE, protocal.MSG_TYPE_RESPONSE, impacket.GetMessageNumber(), nil)
	user.AppendMessage(userResponse)

	// 广播消息
	pushBody := util.JsonMap{"content": content, "user": user.Info}
	targetPush := response.GetSuccessData(config.MESSAGE_ID_BROADCAST_MESSAGE_PUSH, protocal.MSG_TYPE_PUSH, impacket.GetMessageNumber(), pushBody)
	hall.SendBroadcastMessage(targetPush)

	return nil
}

// 消息内容合法性检测、非法字符过滤
func checkAndFilterContent(content string) bool {
	return true
}
