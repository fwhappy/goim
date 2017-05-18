package user

import (
	"sync"

	"github.com/fwhappy/protocal"
	"github.com/fwhappy/util"
)

// User 用户
type User struct {
	Id            string
	Nickname      string
	Avatar        string
	HandshakeTime int64 // 握手时间
	HeartBeatTime int64 // 心跳时间
	Extra         *util.JsonMap
	info          *Info                   // 附加信息
	Mux           *sync.RWMutex           // 用户锁
	Mq            chan *protocal.ImPacket // 消息队列
	RoomIds       *[]string
}

// 创建一个新用户
func NewUser(Id string) *User {
	user := &User{}
	user.Id = Id
	user.Mq = make(chan *protocal.ImPacket, 1024)
	return user
}

// SendMessage 监听用户消息队列，依次给用户发消息
func (user *User) SendMessage() {
	// 捕获异常
	defer util.CatchPanic()

	for imPacket := range user.Mq {
		user.WriteMessage(imPacket)
	}
}
