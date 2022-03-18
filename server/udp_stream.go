package server

import (
	"log"
	"net"
	"strconv"
)

type UdpRemoteStream struct {
	id   uint32
	t    *WebSocketServerTunnel
	conn *net.UDPConn
}

func (s *UdpRemoteStream) Connect(remoteIp string, remotePort uint32) error {
	addr := remoteIp + ":" + strconv.Itoa(int(remotePort))
	raddr, _ := net.ResolveUDPAddr("udp", addr)
	con, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Panicln("connect remote error", err)
		return err
	}
	s.conn = con
	return nil
}

func (s *UdpRemoteStream) WriteToRemote(message []byte) (int, error) {
	return s.conn.Write(message)
}

func (s *UdpRemoteStream) Close() {

	s.conn.Close()
}

func (s *UdpRemoteStream) StartReadRemote() {
	go func() {
		defer s.t.CloseStream(s.id)
		var buf = make([]byte, 8000)
		for {
			size, err := s.conn.Read(buf)
			if err != nil {
				log.Println("stream read remote error:", err)
				return
			}
			log.Printf("recv from remote  msg length = %d", size)
			s.t.Write(s.id, buf[0:size])
		}
	}()
}
