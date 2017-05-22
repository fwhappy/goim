package response

import (
	"goim/core/msgpack"
	"goim/ierror"

	"github.com/fwhappy/protocal"
	"github.com/fwhappy/util"
)

// GetError 返回通用错误 impacket
func GetError(packageType uint8, err *ierror.Error, response util.JsonMap) *protocal.ImPacket {
	if response == nil {
		response = util.JsonMap{}
	}
	response["s2c_result"] = util.JsonMap{"code": err.GetCode(), "msg": err.Error()}
	return protocal.NewImPacket(packageType, msgpack.Marshal(response))
}

// GetSuccess 通用正确消息
func GetSuccess(packageType uint8, response util.JsonMap) *protocal.ImPacket {
	if response == nil {
		response = util.JsonMap{}
	}
	response["s2c_result"] = util.JsonMap{"code": 0}
	return protocal.NewImPacket(packageType, msgpack.Marshal(response))
}
