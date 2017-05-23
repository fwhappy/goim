package room

import (
	"goim/user"
	"sync"
)

// Room 房间信息
type Room struct {
	Id    string // 房间标志
	Mux   *sync.RWMutex
	Users map[string]*user.Info // 用户信息
}

// NewRoom 生成一个新房间
func NewRoom(id string) *Room {
	r := &Room{Id: id, Mux: &sync.RWMutex{}}
	r.Users = make(map[string]*user.Info)
	return r
}

// AddUser 给房间添加一个用户
func (r *Room) AddUser(u *user.User) {
	r.Mux.RLock()
	defer r.Mux.RUnlock()

	r.Users[u.Id] = u.Info
}

// DelUser 删除一个房间用户
func (r *Room) DelUser(id string) {
	r.Mux.RLock()
	defer r.Mux.RUnlock()

	delete(r.Users, id)
}

// IsUserInRoom 判断用户是否在房间内
func (r *Room) IsUserInRoom(id string) bool {
	r.Mux.RLock()
	defer r.Mux.RUnlock()

	_, exists := r.Users[id]
	return exists
}
