package user

import "github.com/fwhappy/util"

// Info 用户共享信息
type Info struct {
	Id       string // 用户id，做一下冗余，方便使用
	Nickname string
	Avatar   string
	Extra    *util.JsonMap
}

// NewInfo 创建一个新的用户共享信息
func NewInfo(id string) *Info {
	return &Info{Id: id}
}
