/*
包内容： 类型 + 消息长度 + 消息
消息 = 消息ID + 消息类型 + 消息编号 + 消息正文
*/

package protocal

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// 各组成部分的长度定义
const (
	PACKAGE_SIZE        = 1 // 包类型
	LENGTH_SIZE         = 3 // 消息长度
	HEADER_SIZE         = 4 // 头，PACKAGE_SIZE + LENGTH_SIZE
	MESSAGE_ID_SIZE     = 2 // 消息ID
	MESSAGE_TYPE_SIZE   = 2 // 消息类型
	MESSAGE_NUMBER_SIZE = 4 // 消息号码
)

// 包类型定义
const (
	PACKAGE_TYPE_HANDSHAKE     = uint8(1)   // 握手
	PACKAGE_TYPE_HANDSHAKE_ACK = uint8(2)   // 握手回复
	PACKAGE_TYPE_HEARTBEAT     = uint8(3)   // 心跳
	PACKAGE_TYPE_DATA          = uint8(4)   // 数据包
	PACKAGE_TYPE_KICK          = uint8(5)   // 退出、踢出
	PACKAGE_TYPE_SYSTEM        = uint8(100) // 系统消息
)

// 消息类型
const (
	MSG_TYPE_REQUEST  = uint16(0)
	MSG_TYPE_NOTIFY   = uint16(1)
	MSG_TYPE_RESPONSE = uint16(2)
	MSG_TYPE_PUSH     = uint16(3)
)

// ImPacket 包
type ImPacket struct {
	buff []byte
}

// Serialize 包序列化成二进制流
func (packet *ImPacket) Serialize() []byte {
	return packet.buff
}

// NewImMessage 生成一个消息
// 消息 = 消息id + 消息类型 + 消息编号 + 消息正文
func NewImMessage(mId uint16, mType uint16, mNumber uint32, body []byte) []byte {
	mBuff := make([]byte, MESSAGE_ID_SIZE+MESSAGE_TYPE_SIZE+MESSAGE_NUMBER_SIZE+len(body))
	// 写入messageId
	binary.BigEndian.PutUint16(mBuff[0:MESSAGE_ID_SIZE], mId)
	// 写入messageType
	binary.BigEndian.PutUint16(mBuff[MESSAGE_ID_SIZE:MESSAGE_ID_SIZE+MESSAGE_TYPE_SIZE], mType)
	// 写入messageNumber
	binary.BigEndian.PutUint32(mBuff[MESSAGE_ID_SIZE+MESSAGE_TYPE_SIZE:MESSAGE_ID_SIZE+MESSAGE_TYPE_SIZE+MESSAGE_NUMBER_SIZE], mNumber)
	// 写入body
	copy(mBuff[MESSAGE_ID_SIZE+MESSAGE_TYPE_SIZE+MESSAGE_NUMBER_SIZE:], body)
	return mBuff
}

// NewImPacket 生成一条消息
func NewImPacket(packageType uint8, message []byte) *ImPacket {
	p := &ImPacket{}
	p.buff = make([]byte, PACKAGE_SIZE+LENGTH_SIZE+len(message))
	// 写入packageType
	p.buff[0] = byte(packageType)
	// 写入包长
	putLength(p.buff[PACKAGE_SIZE:PACKAGE_SIZE+LENGTH_SIZE], uint32(len(message)))
	// 写入包内容
	copy(p.buff[PACKAGE_SIZE+LENGTH_SIZE:], message)
	return p
}

// GetPackage 从字节流中读出包类型
// 兼容老项目，GetPackageType的副函数
func (packet *ImPacket) GetPackage() uint8 {
	return packet.GetPackageType()
}

// GetPackageType 从字节流中读出包类型
func (packet *ImPacket) GetPackageType() uint8 {
	return uint8(packet.buff[0])
}

// GetLength message的长度
// 兼容老项目，GetMessageLength的副函数
func (packet *ImPacket) GetLength() uint32 {
	return packet.GetLength()
}

// GetMessageLength message的长度
func (packet *ImPacket) GetMessageLength() uint32 {
	return length(packet.buff[PACKAGE_SIZE : PACKAGE_SIZE+LENGTH_SIZE])
}

// GetMessage 读取消息内容
func (packet *ImPacket) GetMessage() []byte {
	return packet.buff[HEADER_SIZE:]
}

// GetMessageId 解析数据包的消息id，非数据包，直接返回0
func (packet *ImPacket) GetMessageId() uint16 {
	if packet.GetPackageType() == PACKAGE_TYPE_DATA {
		message := packet.GetMessage()
		return binary.BigEndian.Uint16(message[0:MESSAGE_ID_SIZE])
	}
	return uint16(0)
}

// GetMessageType 解析数据包的消息类型，非数据包，直接返回0
func (packet *ImPacket) GetMessageType() uint16 {
	if packet.GetPackageType() == PACKAGE_TYPE_DATA {
		message := packet.GetMessage()
		return binary.BigEndian.Uint16(message[MESSAGE_ID_SIZE : MESSAGE_ID_SIZE+MESSAGE_TYPE_SIZE])
	}
	return uint16(0)
}

// GetMessageNumber 解析数据包的消息编号，非数据包，直接返回0
func (packet *ImPacket) GetMessageNumber() uint32 {
	if packet.GetPackageType() == PACKAGE_TYPE_DATA {
		message := packet.GetMessage()
		return binary.BigEndian.Uint32(message[MESSAGE_ID_SIZE+MESSAGE_TYPE_SIZE : MESSAGE_ID_SIZE+MESSAGE_TYPE_SIZE+MESSAGE_NUMBER_SIZE])
	}
	return uint32(0)
}

// GetBody 解析数据包的消息正文
func (packet *ImPacket) GetBody() []byte {
	message := packet.GetMessage()
	if packet.GetPackageType() == PACKAGE_TYPE_DATA {
		return message[MESSAGE_ID_SIZE+MESSAGE_TYPE_SIZE+MESSAGE_NUMBER_SIZE:]
	}
	return message
}

// ReadPacket 从socket中读出一条小
func ReadPacket(conn *net.TCPConn) (*ImPacket, error) {
	var (
		packageBytes = make([]byte, PACKAGE_SIZE)
		lengthBytes  = make([]byte, LENGTH_SIZE)
		packageType  uint8
	)

	// 读取package
	if _, err := io.ReadFull(conn, packageBytes); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, err
		}
		return nil, fmt.Errorf("packageType read error: %s", err.Error())
	}
	// 转成uint8
	packageType = packageBytes[0]

	// 读取lengthBytes
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, err
		}
		return nil, fmt.Errorf("packet length read error: %s", err.Error())
	}
	// 内容长度
	mLength := length(lengthBytes)

	// 读取message
	messageBytes := make([]byte, mLength)
	if _, err := io.ReadFull(conn, messageBytes); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, err
		}
		return nil, fmt.Errorf("read packet message error: %s", err.Error())
	}
	return NewImPacket(packageType, messageBytes), nil
}

// 写入长度
func putLength(b []byte, v uint32) {
	_ = b[2] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

// 获取长度
func length(b []byte) uint32 {
	_ = b[2] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
}

// Send 发送消息至socket连接
func (packet *ImPacket) Send(conn *net.TCPConn) {
	conn.Write(packet.Serialize())
}
