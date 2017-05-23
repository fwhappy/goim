package room

import (
	"goim/config"
	"goim/core"
	"goim/core/msgpack"
	"goim/ierror"
	"goim/response"
	"net"

	"github.com/fwhappy/protocal"

	roomService "goim/service/room"
)

// JoinRoom 加入房间
func JoinRoom(id string, conn *net.TCPConn, impacket *protocal.ImPacket) {
	var err *ierror.Error
	body, _ := msgpack.Unmarshal(impacket.GetBody())
	roomId, _ := body.JsonGetString("room_id")
	if roomId == "" {
		err = ierror.NewError(-301, "rooom_id")
	} else {
		err = roomService.JoinRoom(id, roomId, impacket)
	}

	if err != nil {
		core.Logger.Error("[JoinRoom]%v", err.Error())
		response.GetErrorData(config.MESSAGE_ID_JOIN_ROOM_RESPONSE, protocal.MSG_TYPE_RESPONSE, impacket.GetMessageNumber(), err, nil).Send(conn)
	}
}

// QuitRoom 退出房间
func QuitRoom(id string, conn *net.TCPConn, impacket *protocal.ImPacket) {
	var err *ierror.Error
	body, _ := msgpack.Unmarshal(impacket.GetBody())
	roomId, _ := body.JsonGetString("room_id")
	if roomId == "" {
		err = ierror.NewError(-302, "rooom_id")
	} else {
		err = roomService.QuitRoom(id, roomId, impacket)
	}

	if err != nil {
		core.Logger.Error("[QuitRoom]%v", err.Error())
		response.GetErrorData(config.MESSAGE_ID_QUIT_ROOM_RESPONSE, protocal.MSG_TYPE_RESPONSE, impacket.GetMessageNumber(), err, nil).Send(conn)
	}
}
