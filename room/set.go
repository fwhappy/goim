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

// Get 根据Id读取用户信息
// 这里不加锁，需要在业务层手动加锁
func (s *Set) Get(id string) (*Room, bool) {
	r, exists := s.Rooms[id]
	return r, exists
}

// Add 将用户加入到用户集合
// 这里不加锁，需要在业务层手动加锁
func (s *Set) Add(r *Room) {
	s.Rooms[r.Id] = r
}

// Del 将用户从集合移除
// 这里不加锁，需要在业务层手动加锁
func (s *Set) Del(id string) {
	delete(s.Rooms, id)
}

// IsExists 判断用户是否存在于集合中
// 这里不加锁，需要在业务层手动加锁
func (s *Set) IsExists(id string) bool {
	_, exists := s.Rooms[id]
	return exists
}
