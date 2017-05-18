package room

import (
	"sync"
)

// Room 房间信息
type Room struct {
	ID    string // 房间标志
	Mux   *sync.RWMutex
	Users []string // 房间用户id列表
}

// NewRoom 生成一个新房间
func NewRoom(id string) (*Room, error) {
	r := &Room{ID: id, Mux: &sync.RWMutex{}}
	return r, nil
}

// AddUser 给房间添加一个用户
func (r *Room) AddUser(u *User) error {
	r.Mux.RLock()
	defer r.Mux.RUnlock()

	r.Users = append(r.Users, u.UserId)
	return nil
}

// Len 房间用户列表
func (r *Room) Len() int {
	return len(r.Users)
}