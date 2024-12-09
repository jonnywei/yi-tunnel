package client

import (
	"net"
	"yi_tunnel/common"
)

type ITunnel interface {
	Open() error
	CreateStream(local net.Conn) *Stream
	Write(stream *Stream, message []byte)
	CloseStream(stream *Stream)
	IsClosed() bool
	StreamCount() int
}

var TunnelIdGen uint32

type BaseTunnel struct {
	UdpStreams map[string]*Stream
	Streams    map[uint32]*Stream
	Config     *common.Config
	Name       string
}
