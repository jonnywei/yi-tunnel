package client

import (
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jonnywei/yi_tunnel/common"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type WebSocketTunnel struct {
	common.CommandBuilder
	BaseTunnel
	mux    sync.Mutex
	smux   sync.RWMutex
	conn   *websocket.Conn
	closed bool
}

func NewWebSocketTunnel(config *common.Config) *WebSocketTunnel {

	tunnel := WebSocketTunnel{
		BaseTunnel: BaseTunnel{
			Config:     config,
			Streams:    make(map[uint32]*Stream),
			UdpStreams: make(map[string]*Stream),
			Name:       fmt.Sprintf("tunnel%d", atomic.AddUint32(&TunnelIdGen, 1)),
		},
	}
	return &tunnel
}

func (t *WebSocketTunnel) CreateStream(local net.Conn) *Stream {
	stream := NewStream(t, local)
	t.smux.Lock()
	t.Streams[stream.Id] = stream
	t.smux.Unlock()
	buf := t.BuildOpenCommand(stream.Id)
	log.Printf("client new stream %d ", stream.Id)
	t.sendBytes(buf)
	return stream
}

func (t *WebSocketTunnel) CreateUdpStreamOrSend(u *UdpLocal, laddr *net.UDPAddr, message []byte) *Stream {
	sa := laddr.String()
	stream, ok := t.UdpStreams[sa]
	if !ok {
		stream = t.CreateUdpStream(u, laddr, message)
		t.UdpStreams[sa] = stream
	} else {
		t.Write(stream, message)
	}
	return stream
}

func (t *WebSocketTunnel) CreateUdpStream(u *UdpLocal, laddr *net.UDPAddr, message []byte) *Stream {
	stream := NewUdpStream(t, laddr, u)
	t.smux.Lock()
	t.Streams[stream.Id] = stream
	t.smux.Unlock()
	buf := t.BuildUdpOpenCommand(stream.Id, message)
	log.Printf("new udp stream %d ", stream.Id)
	t.sendBytes(buf)
	return stream
}

func (t *WebSocketTunnel) Write(stream *Stream, message []byte) {
	buf := t.BuildDataCommand(stream.Id, message)
	log.Printf("client stream %d write data to tunnel len=%d ", stream.Id, len(buf))
	t.sendBytes(buf)
	log.Printf("client stream %d write data to tunnel len=%d  end", stream.Id, len(buf))

}

func (t *WebSocketTunnel) OnOpen() {
	buf := t.BuildInitCommand(t.Config.Remote_address, t.Config.Remote_port, t.Config.Secret)
	t.sendBytes(buf)
}

func (t *WebSocketTunnel) sendBytes(message []byte) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.conn.WriteMessage(websocket.BinaryMessage, message)
}

func (t *WebSocketTunnel) StreamCount() int {
	t.smux.RLock()
	defer t.smux.RUnlock()
	return len(t.Streams)
}

func (t *WebSocketTunnel) OnMessage(message []byte) {

	length := binary.BigEndian.Uint32(message)
	command := binary.BigEndian.Uint32(message[4:8])
	streamId := binary.BigEndian.Uint32(message[8:12])
	log.Printf("client receive length:%d,command:%d,streamId:%d\n", length, command, streamId)
	t.smux.RLock()
	stream := t.Streams[streamId]
	t.smux.RUnlock()
	if stream == nil {
		log.Printf("client stream %d cannot find\n", streamId)
		return
	}
	if command == common.CommandOpen {
		stream.Open()
	} else if command == common.CommandUdpOpen {
		//stream.Open()
	} else if command == common.CommandData {
		stream.WriteToLocal(message[12:])
	} else if command == common.CommandClose {
		t.CloseStream(stream)
	}

}

func (t *WebSocketTunnel) Open() error {

	var wsurl = t.Config.Tunnel_config
	c, _, err := websocket.DefaultDialer.Dial(wsurl, nil)
	if err != nil {
		log.Printf("connect tunnel error %s\n", wsurl)
		log.Println(err)
		return err
	}
	t.conn = c
	t.OnOpen()
	go func() {
		defer t.OnClose()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("web socket client recv from tunnel  msg length = %d", len(message))
			t.OnMessage(message)
		}
	}()
	return nil
}

func (t *WebSocketTunnel) CloseStream(stream *Stream) {

	log.Printf("stream %d send close command to server\n", stream.Id)
	buffer := t.BuildCloseCommand(stream.Id)
	t.smux.Lock()
	delete(t.Streams, stream.Id)
	t.smux.Unlock()
	t.sendBytes(buffer)
	stream.Close()
}

func (t *WebSocketTunnel) OnClose() {
	log.Print("client close websocket tunnel")
	//
	t.closed = true
	t.smux.RLock()
	for _, stream := range t.Streams {
		stream.Close()
	}
	t.smux.RUnlock()

}

func (t *WebSocketTunnel) IsClosed() bool {
	return t.closed
}
