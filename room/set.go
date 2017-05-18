package room

import "sync"

// Set 房间集合
type Set struct {
	Mux   *sync.RWMutex
	Rooms map[string]*Room
}

// NewSet 生成一个房间列表
func NewSet() *Set {
	s := &Set{}
	s.Rooms = make(map[string]*Room)
	s.Mux = &sync.RWMutex{}
	return s
}

// Len 房间个数
func (s *Set) Len() int {
	return len(s.Rooms)
}
