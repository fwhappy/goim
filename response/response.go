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
	response["s2c_result"] = util.JsonMap{"code": 0, "msg": ""}
	return protocal.NewImPacket(packageType, msgpack.Marshal(response))
}

// GetErrorData 返回通用数据类错误 impacket
func GetErrorData(messageId, messageType uint16, messageNumber uint32, err *ierror.Error, response util.JsonMap) *protocal.ImPacket {
	if response == nil {
		response = util.JsonMap{}
	}
	response["s2c_result"] = util.JsonMap{"code": err.GetCode(), "msg": err.Error()}
	message := protocal.NewImMessage(messageId, messageType, messageNumber, msgpack.Marshal(response))
	return protocal.NewImPacket(protocal.PACKAGE_TYPE_DATA, message)
}

// GetSuccessData 通用数据回应
func GetSuccessData(messageId, messageType uint16, messageNumber uint32, response util.JsonMap) *protocal.ImPacket {
	if response == nil {
		response = util.JsonMap{}
	}
	response["s2c_result"] = util.JsonMap{"code": 0, "msg": ""}
	message := protocal.NewImMessage(messageId, messageType, messageNumber, msgpack.Marshal(response))
	return protocal.NewImPacket(protocal.PACKAGE_TYPE_DATA, message)
}
