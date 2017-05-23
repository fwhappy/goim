package hall

import (
	"goim/config"
	"goim/core"
	"goim/response"
	"goim/room"
	"goim/user"

	"github.com/fwhappy/protocal"
	"github.com/fwhappy/util"
)

// 定义全局变量
var (
	UserSet *user.Set
	RoomSet *room.Set
)

func init() {
	UserSet = user.NewSet()
	RoomSet = room.NewSet()
}

// SendUserMessage 通过Id给用户发送一条消息
func SendUserMessage(id string, imPacket *protocal.ImPacket) bool {
	if u, online := UserSet.Get(id); online {
		u.AppendMessage(imPacket)
		return true
	}
	return false
}

// SendRoomMessage 推送消息给房间用户
func SendRoomMessage(r *room.Room, impacket *protocal.ImPacket) {
	if len(r.Users) == 0 {
		return
	}

	r.Mux.Lock()
	defer r.Mux.Unlock()

	for id := range r.Users {
		SendUserMessage(id, impacket)
	}
}

// SendBroadcastMessage 推送广播消息
func SendBroadcastMessage(impacket *protocal.ImPacket) {
	if UserSet.Len() == 0 {
		return
	}

	// 这里是全局用户锁，会大大影响性能，肯定不能这么搞，需要后续的解决方案
	UserSet.Mux.Lock()
	defer UserSet.Mux.Unlock()

	for _, u := range UserSet.Users {
		u.AppendMessage(impacket)
	}
}

// RemoveRoomUser 移除房间成员
func RemoveRoomUser(roomId string, u *user.User) {
	// 读取房间信息
	if r, exists := RoomSet.Get(roomId); exists {
		// 从成员列表删除
		r.DelUser(u.Id)
		core.Logger.Debug("[QuitRoom]roomId:%v, id:%v", r.Id, u.Id)

		if len(r.Users) == 0 {
			// 房间没有其他用户，解散房间
			RoomSet.Del(r.Id)
			core.Logger.Debug("[DismissRoom]roomId:%v", r.Id)
		} else {
			// 通知其他用户，有人退出了房间
			body := make(util.JsonMap)
			body["user"] = u.Info
			pushPacket := response.GetSuccessData(config.MESSAGE_ID_QUIT_ROOM_PUSH, protocal.MSG_TYPE_PUSH, 0, body)
			SendRoomMessage(r, pushPacket)
		}
	}
}
