package common

import "encoding/binary"

const (
	CommandOpen    = 1
	CommandData    = 2
	CommandClose   = 3
	CommandInit    = 4
	CommandUdpOpen = 5
	DefaultLength  = 0
)

type CommandBuilder struct {
}

func (t *CommandBuilder) BuildOpenCommand(streamId uint32) []byte {

	buf := make([]byte, 12)
	binary.BigEndian.PutUint32(buf, 12)
	binary.BigEndian.PutUint32(buf[4:], CommandOpen)
	binary.BigEndian.PutUint32(buf[8:], streamId)
	return buf
}

func (t *CommandBuilder) BuildInitCommand(remoteAddress string, remotePort int, secret string) []byte {

	buf := make([]byte, 12)
	binary.BigEndian.PutUint32(buf, DefaultLength)
	binary.BigEndian.PutUint32(buf[4:], CommandInit)
	secretb := []byte(secret)
	binary.BigEndian.PutUint32(buf[8:], uint32(len(secretb)))
	buf = append(buf, secretb...)
	addb := []byte(remoteAddress)
	temp4 := make([]byte, 4)
	binary.BigEndian.PutUint32(temp4, uint32(len(addb)))
	buf = append(buf, temp4...)
	buf = append(buf, addb...)
	temp4 = make([]byte, 4)
	binary.BigEndian.PutUint32(temp4, uint32(remotePort))
	buf = append(buf, temp4...)
	binary.BigEndian.PutUint32(buf[0:], uint32(len(buf)))
	return buf

}

func (t *CommandBuilder) BuildUdpOpenCommand(streamId uint32, message []byte) []byte {

	buf := make([]byte, 12)
	binary.BigEndian.PutUint32(buf, DefaultLength)
	binary.BigEndian.PutUint32(buf[4:], CommandUdpOpen)
	binary.BigEndian.PutUint32(buf[8:], streamId)
	buf = append(buf, message...)
	binary.BigEndian.PutUint32(buf[0:], uint32(len(buf)))
	return buf
}

func (t *CommandBuilder) BuildDataCommand(streamId uint32, message []byte) []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint32(buf, DefaultLength)
	binary.BigEndian.PutUint32(buf[4:], CommandData)
	binary.BigEndian.PutUint32(buf[8:], streamId)
	buf = append(buf, message...)
	binary.BigEndian.PutUint32(buf[0:], uint32(len(buf)))
	return buf
}

func (t *CommandBuilder) BuildCloseCommand(streamId uint32) []byte {

	buf := make([]byte, 12)
	binary.BigEndian.PutUint32(buf, 12)
	binary.BigEndian.PutUint32(buf[4:], CommandClose)
	binary.BigEndian.PutUint32(buf[8:], streamId)
	return buf
}
