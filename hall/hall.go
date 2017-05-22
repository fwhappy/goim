package hall

import (
	"goim/room"
	"goim/user"
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
