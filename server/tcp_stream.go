package server

import (
	"log"
	"net"
	"strconv"
)

type RemoteStream interface {
	Connect(remoteIp string, remotePort uint32) error

	Close()

	WriteToRemote(message []byte) (int, error)

	StartReadRemote()
}

type TcpRemoteStream struct {
	id   uint32
	t    *WebSocketServerTunnel
	conn *net.TCPConn
}

func (s *TcpRemoteStream) Connect(remoteIp string, remotePort uint32) error {
	addr := remoteIp + ":" + strconv.Itoa(int(remotePort))
	raddr, _ := net.ResolveTCPAddr("tcp", addr)
	con, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		log.Println("connect remote error", err)
		return err
	}
	s.conn = con
	return nil
}

func (s *TcpRemoteStream) WriteToRemote(message []byte) (int, error) {
	return s.conn.Write(message)
}

func (s *TcpRemoteStream) Close() {

	s.conn.Close()
}

func (s *TcpRemoteStream) StartReadRemote() {
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
