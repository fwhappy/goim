package user

import "sync"

// Set 用户集合
type Set struct {
	Mux   *sync.RWMutex
	Users map[string]*User
}

// NewSet 生成一个用户集合
func NewSet() *Set {
	s := &Set{}
	s.Users = make(map[string]*User)
	s.Mux = &sync.RWMutex{}
	return s
}

// Len 用户个数
func (s *Set) Len() int {
	return len(s.Users)
}

// Get 根据Id读取用户信息
func (s *Set) Get(id string) (*User, bool) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	u, exists := s.Users[id]
	return u, exists
}

// Add 将用户加入到用户集合
func (s *Set) Add(u *User) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	s.Users[u.Id] = u
}

// Del 将用户从集合移除
func (s *Set) Del(id string) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	delete(s.Users, id)
}

// IsExists 判断用户是否存在于集合中
func (s *Set) IsExists(Id string) bool {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	_, exists := s.Users[Id]
	return exists
}
