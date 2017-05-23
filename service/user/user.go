package user

import (
	"goim/core"
	"goim/core/msgpack"
	"goim/hall"
	"goim/ierror"
	"goim/response"
	"goim/user"
	"net"
	"time"

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
	u.HandshakeTime = util.GetTime()
	u.Info.Nickname = nickname
	u.Info.Extra = &extra
	hall.UserSet.Add(u)

	core.Logger.Debug("[HandShake]id:%v, user:%v", id, u)

	// 通知用户握手成功
	body := util.JsonMap{"heartbeat_interval": core.GetAppConfig("heartbeat_interval").(int64)}
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

	// 开启心跳检测
	go listenHeartBeat(u)

	core.Logger.Info("[HandShakeAck]id:%v", id)

	return nil
}

// HeartBeat 用户心跳
func HeartBeat(id string) {
	if u, online := hall.UserSet.Get(id); online {
		u.HeartBeatTime = util.GetTime()
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

// 监听用户心跳
func listenHeartBeat(u *user.User) {
	// 捕获异常
	defer util.RecoverPanic()

	// 读取心跳间隔
	heartBeatInterval := core.GetAppConfig("heartbeat_interval").(int64)

	for {
		time.Sleep(time.Second * time.Duration(heartBeatInterval))
		user, online := hall.UserSet.Get(u.Id)
		if !online {
			core.Logger.Debug("用户已下线，停止心跳监测, id:%v", u.Id)
			break
		}

		if user.HandshakeTime != u.HandshakeTime {
			// 用户已经被顶号或者重新登录
			core.Logger.Debug("用户已重新登录，停止心跳监测, id:%v", u.Id)
			break
		}

		if util.GetTime()-user.HeartBeatTime > 2*heartBeatInterval {
			core.Logger.Debug("用户心跳停止，踢下线, id:%v", u.Id)
			KickUser(u)
			break
		}
	}
}
