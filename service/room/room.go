package room

import (
	"goim/config"
	"goim/core"
	"goim/hall"
	"goim/ierror"
	"goim/response"
	"goim/room"

	"github.com/fwhappy/protocal"
	"github.com/fwhappy/util"
)

// JoinRoom 加入房间
func JoinRoom(id string, roomId string, impacket *protocal.ImPacket) *ierror.Error {
	// 读取用户信息
	u, online := hall.UserSet.Get(id)
	if !online {
		return ierror.NewError(-202)
	}

	hall.RoomSet.Mux.Lock()
	defer hall.RoomSet.Mux.Unlock()

	// 创建一个新的房间
	r, exists := hall.RoomSet.Get(roomId)
	if !exists {
		r = room.NewRoom(roomId)
		hall.RoomSet.Add(r)
	}

	// 返回一个回应，通知加入房间成功
	u.AppendMessage(response.GetSuccessData(config.MESSAGE_ID_JOIN_ROOM_RESPONSE, protocal.MSG_TYPE_RESPONSE, impacket.GetMessageNumber(), nil))

	if !r.IsUserInRoom(id) {
		// 通知其他用户，有人加入了房间
		body := make(util.JsonMap)
		body["user"] = u.Info
		pushPacket := response.GetSuccessData(config.MESSAGE_ID_JOIN_ROOM_PUSH, protocal.MSG_TYPE_PUSH, 0, body)
		hall.SendRoomMessage(r, pushPacket)
		// 将用户加入房间
		r.AddUser(u)
		u.AddRoom(r.Id)
		core.Logger.Info("[JoinRoom]roomId:%v,id:%v", r.Id, id)
	}

	return nil
}

// QuitRoom 退出房间
func QuitRoom(id string, roomId string, impacket *protocal.ImPacket) *ierror.Error {
	// 读取用户信息
	u, online := hall.UserSet.Get(id)
	if !online {
		return ierror.NewError(-202)
	}

	// 删除房间用户
	hall.RemoveRoomUser(roomId, u)

	// 从用户的房间列表中删除房间id
	u.DelRoom(roomId)

	// 返回一个回应，通知退出房间成功
	u.AppendMessage(response.GetSuccessData(config.MESSAGE_ID_QUIT_ROOM_RESPONSE, protocal.MSG_TYPE_RESPONSE, impacket.GetMessageNumber(), nil))

	return nil
}
