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

// SendMessage 通过Id给用户发送一条消息
func SendMessage(id string, imPacket *protocal.ImPacket) bool {
	if user, online := hall.UserSet.Get(id); online {
		user.AppendMessage(imPacket)
		return true
	}
	return false
}

// HandShake 用户握手
func HandShake(conn *net.TCPConn, impacket *protocal.ImPacket) (string, *ierror.Error) {
	// 解析用户数据
	info, err := msgpack.Unmarshal(impacket.GetBody())
	if err != nil {
		core.Logger.Debug("decode error:%v", err.Error())
		return "", ierror.NewError(-200)
	}

	// check数据完整性
	id, exists := info.JsonGetString("id")
	if !exists {
		return "", ierror.NewError(-201, "id")
	}
	nickname, exists := info.JsonGetString("nickname")
	if !exists {
		return "", ierror.NewError(-201, "nickname")
	}
	// 附加信息
	extra := info.JsonGetJsonMap("extra")

	// 判断用户是否已在线
	if u, online := hall.UserSet.Get("id"); online {
		// TODO 踢老的连接下线
		// TODO 通知用户下线
		u.Conn.Close()

		core.Logger.Debug("[repeat handshake], id:%v, old remote:%v, new remote:%v", id, u.Conn.RemoteAddr().String(), conn.RemoteAddr().String())
	}

	u := user.NewUser(id)
	u.Conn = conn
	u.Nickname = nickname
	u.Extra = &extra
	hall.UserSet.Add(u)

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

// 踢出用户
func KickUser(u *user.User) {
	u.QuitOnce.Do(func() {
		// 将用户从大厅删除
		hall.UserSet.Del(u.Id)

		// 关闭用户连接
		u.Conn.Close()

		// 关闭用户消息队列
		close(u.Mq)

		// TODO 将用户从房间移除

		core.Logger.Info("[kickUser]id:%v", u.Id)
	})
}
