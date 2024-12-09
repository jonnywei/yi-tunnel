package server

import (
	"encoding/binary"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/jonnywei/yi_tunnel/common"
	"log"
	"sync"
)

type WebSocketServerTunnel struct {
	sync.Mutex
	common.CommandBuilder
	Streams       map[uint32]RemoteStream
	remoteAddress string
	remotePort    uint32
	con           *websocket.Conn
	secret        string
}

func NewWebSocketServerTunnel(con *websocket.Conn, secret string) *WebSocketServerTunnel {

	return &WebSocketServerTunnel{
		Streams: make(map[uint32]RemoteStream),
		con:     con,
		secret:  secret,
	}
}

func (t *WebSocketServerTunnel) OnMessage(message []byte) error {

	length := binary.BigEndian.Uint32(message)
	command := binary.BigEndian.Uint32(message[4:8])
	if command == common.CommandInit {
		log.Printf("server receive length:%d,command:Init", length)
		secretlen := binary.BigEndian.Uint32(message[8:12])
		clientSecret := string(message[12 : 12+secretlen])
		if clientSecret != t.secret {
			return errors.New("error secret key")
		}
		index := 12 + secretlen
		addrLen := binary.BigEndian.Uint32(message[index : index+4])
		index = 4 + index
		t.remoteAddress = string(message[index : index+addrLen])
		index = addrLen + index
		t.remotePort = binary.BigEndian.Uint32(message[index:])
		return nil
	}
	streamId := binary.BigEndian.Uint32(message[8:12])
	log.Printf("server receive length:%d,command:%d,streamId:%d\n", length, command, streamId)

	if command == common.CommandOpen {
		stream := TcpRemoteStream{t: t, id: streamId}
		log.Println("connect to remote", t.remoteAddress, t.remotePort)
		err := stream.Connect(t.remoteAddress, t.remotePort)
		if err != nil {
			log.Println("connect error", err)
			buf := t.BuildCloseCommand(streamId)
			t.writeToTunnel(buf)
		}
		t.Streams[streamId] = &stream
		buf := t.BuildOpenCommand(streamId)
		t.writeToTunnel(buf)
		stream.StartReadRemote()
	}

	if command == common.CommandUdpOpen {
		stream := UdpRemoteStream{t: t, id: streamId}
		log.Println("connect to remote", t.remoteAddress, t.remotePort)
		err := stream.Connect(t.remoteAddress, t.remotePort)
		if err != nil {
			log.Println("connect error", err)
			buf := t.BuildCloseCommand(streamId)
			t.writeToTunnel(buf)
		}
		t.Streams[streamId] = &stream
		buf := t.BuildUdpOpenCommand(streamId, []byte{0})
		t.writeToTunnel(buf)
		stream.StartReadRemote()
		//xieshuju
		stream.WriteToRemote(message[12:])
	}
	stream := t.Streams[streamId]
	if stream == nil {
		log.Printf("server stream %d cannot find\n", streamId)
		return nil
	}
	if command == common.CommandData {
		stream.WriteToRemote(message[12:])

	} else if command == common.CommandClose {
		log.Printf("stream %d receive peer close command\n", streamId)
		t.CloseStream(streamId)
	}
	return nil
}

func (t *WebSocketServerTunnel) OnClose() {
	for k, _ := range t.Streams {
		t.CloseStream(k)
	}
	t.con.Close()
}

func (t *WebSocketServerTunnel) Write(streamId uint32, message []byte) {
	buf := t.BuildDataCommand(streamId, message)
	log.Printf(" stream %d write data to tunnel len=%d ", streamId, len(buf))
	t.writeToTunnel(buf)
	log.Printf(" stream %d write data to tunnel len=%d end ", streamId, len(buf))

}

func (t *WebSocketServerTunnel) CloseStream(streamId uint32) {
	stream := t.Streams[streamId]
	if stream == nil {
		log.Printf("close stream %d cannot find\n", streamId)
		return
	}
	delete(t.Streams, streamId)
	buf := t.BuildCloseCommand(streamId)
	t.writeToTunnel(buf)
	stream.Close()
	log.Printf("server stream %d closed \n", streamId)
}

func (t *WebSocketServerTunnel) writeToTunnel(buf []byte) {
	t.Lock()
	defer t.Unlock()
	t.con.WriteMessage(websocket.BinaryMessage, buf)
}
