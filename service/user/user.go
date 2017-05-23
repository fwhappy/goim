package user

import (
	"goim/core"
	"goim/core/msgpack"
	"goim/hall"
	"goim/ierror"
	"goim/response"
	"goim/user"
	"net"

	"github.com/fwhappy/protocal"
	"github.com/fwhappy/util"
)

// HandShake 用户握手
func HandShake(conn *net.TCPConn, impacket *protocal.ImPacket) (string, *ierror.Error) {
	// 解析用户数据
	info, err := msgpack.Unmarshal(impacket.GetBody())
	if err != nil {
		return "", ierror.NewError(-100, "HandShake", "")
	}

	// check数据完整性
	id, exists := info.JsonGetString("id")
	if !exists {
		return "", ierror.NewError(-201, id)
	}
	nickname, exists := info.JsonGetString("nickname")
	if !exists {
		return "", ierror.NewError(-201, nickname)
	}
	// 附加信息
	extra := info.JsonGetJsonMap("extra")
	core.Logger.Debug("[HandShake]id:%v, extra:%v", id, extra)

	// 判断用户是否已在线
	if connectedUser, online := hall.UserSet.Get(id); online {
		core.Logger.Debug("[repeat handshake], id:%v, old remote:%v, new remote:%v", id, connectedUser.Conn.RemoteAddr().String(), conn.RemoteAddr().String())
		KickUser(connectedUser)
	}

	u := user.NewUser(id)
	u.Conn = conn
	u.Info.Nickname = nickname
	u.Info.Extra = &extra
	hall.UserSet.Add(u)

	core.Logger.Debug("[HandShake]id:%v, user:%v", id, u)

	// 通知用户握手成功
	body := util.JsonMap{"heartbeat_interval": core.GetAppConfig("heartbeat_interval").(string)}
	u.WriteMessage(response.GetSuccess(protocal.PACKAGE_TYPE_HANDSHAKE, body))

	core.Logger.Info("[HandShake]id:%v", id)

	return u.Id, nil
}

// HandShakeAck 握手成功
func HandShakeAck(id string, impacket *protocal.ImPacket) *ierror.Error {
	u, online := hall.UserSet.Get(id)
	if !online {
		return ierror.NewError(-202)
	}

	// 通知回复成功
	u.WriteMessage(response.GetSuccess(protocal.PACKAGE_TYPE_HANDSHAKE_ACK, nil))

	// 开启消息推送
	go u.SendMessage()

	core.Logger.Info("[HandShakeAck]id:%v", id)

	return nil
}

// HeartBeat 用户心跳
func HeartBeat(id string) {
	if u, online := hall.UserSet.Get(id); online {
		u.WriteMessage(response.GetSuccess(protocal.PACKAGE_TYPE_HEARTBEAT, nil))
		// 适时屏蔽
		// core.Logger.Debug("[HeartBeat]id:%v", id)
	} else {
		core.Logger.Warn("[HeartBeat]user not online,id:%v", id)
	}
}

// Logout 用户请求退出
func Logout(id string) *ierror.Error {
	u, online := hall.UserSet.Get(id)
	if !online {
		return ierror.NewError(-202)
	}

	// 剔除用户
	KickUser(u)

	return nil
}

// KickUser 踢出用户
func KickUser(u *user.User) {
	u.QuitOnce.Do(func() {
		// 将用户从大厅删除
		hall.UserSet.Del(u.Id)

		// 关闭用户连接
		u.Conn.Close()
		core.Logger.Debug("u.Conn=%v", u.Conn.RemoteAddr().String())

		// 关闭用户消息队列
		close(u.Mq)

		// 将用户从房间移除
		for _, roomId := range u.RoomIds {
			// 删除房间用户
			hall.RemoveRoomUser(roomId, u)
		}

		core.Logger.Info("[kickUser]id:%v", u.Id)
	})
}
