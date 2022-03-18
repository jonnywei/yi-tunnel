package client

import (
	"io"
	"log"
	"net"
	"sync/atomic"
)

const (
	StreamTypeTCP StreamType = 0x01
	StreamTypeUDP StreamType = 0x02
)

type StreamType int

type Stream struct {
	sType    StreamType
	tunnel   ITunnel
	local    net.Conn
	addr     *net.UDPAddr
	Id       uint32
	udplocal *UdpLocal
}

var IdGen uint32

func NewStream(tunnel ITunnel, local net.Conn) *Stream {
	return &Stream{
		sType:  StreamTypeTCP,
		tunnel: tunnel,
		local:  local,
		Id:     atomic.AddUint32(&IdGen, 1),
	}
}

func NewUdpStream(tunnel ITunnel, local *net.UDPAddr, u *UdpLocal) *Stream {
	return &Stream{
		sType:    StreamTypeUDP,
		tunnel:   tunnel,
		addr:     local,
		udplocal: u,
		Id:       atomic.AddUint32(&IdGen, 1),
	}
}

func (t *Stream) Open() {
	if t.sType == StreamTypeTCP {
		go func() {
			var localClose = false
			for {
				var buf = make([]byte, 4096)
				n, err := t.local.Read(buf)
				log.Printf("tcp local read %d byte", n)
				if err != nil {
					log.Print("conn read error ", err)
					if err == io.EOF {
						localClose = true
					}
					break
				}
				t.WriteToTunnel(buf[0:n])
			}
			if localClose {
				t.CloseStream()
			}
		}()
	}

}

func (t *Stream) WriteToLocal(message []byte) {
	if t.sType == StreamTypeTCP {
		t.local.Write(message)
	}
	if t.sType == StreamTypeUDP {
		t.udplocal.WriteToLocal(t.addr, message)
	}
}

func (t *Stream) WriteToTunnel(message []byte) {
	t.tunnel.Write(t, message)
}

func (t *Stream) Close() {
	log.Printf("stream %d close", t.Id)
	if t.sType == StreamTypeTCP {
		err := t.local.Close()
		if err != nil {
			log.Println("stream local close error", err)
		}
	}
	if t.sType == StreamTypeUDP {

	}
}
func (t *Stream) CloseStream() {
	t.tunnel.CloseStream(t)

}
