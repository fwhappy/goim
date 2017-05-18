package protocal

import "testing"

// TestPacket 非DATA类型的数据包
func TestPacket(t *testing.T) {
	packageType := PACKAGE_TYPE_HANDSHAKE
	content := "hello, i am a handshake message"
	packet := NewImPacket(packageType, []byte(content))

	if packageType != packet.GetPackageType() {
		t.Error("TestPacket GetPackageType error.")
	}
	if string(packet.GetBody()) != content {
		t.Error("TestPacket GetBody error.")
	}
}

// TestPackData DATA类型的数据包
func TestPackData(t *testing.T) {
	packageType := uint8(4)
	messageId := uint16(1)
	messageType := MSG_TYPE_REQUEST
	messageNumber := uint32(64)
	content := "hello, i am a data message"
	message := NewImMessage(messageId, messageType, messageNumber, []byte(content))
	packet := NewImPacket(packageType, message)

	if packageType != packet.GetPackageType() {
		t.Error("TestPacket GetPackageType error.")
	}
	if messageId != packet.GetMessageId() {
		t.Error("TestPacket GetMessageId error.")
	}
	if messageType != packet.GetMessageType() {
		t.Error("TestPacket GetMessageType error.")
	}
	if messageNumber != packet.GetMessageNumber() {
		t.Error("TestPacket GetMessageNumber error.")
	}
	if string(packet.GetBody()) != content {
		t.Error("TestPacket GetBody error.")
	}
}
