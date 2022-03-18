package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"yi-tunnel/common"
)

type UdpServer struct {
	Config  common.ServerConfig
	udpAddr *net.UDPAddr
	tunnel  WebSocketServerTunnel
}

func NewUdpServer(config *common.ServerConfig) *UdpServer {

	tcpLocal := UdpServer{Config: *config}
	return &tcpLocal
}

func (u *UdpServer) listen() {

	u.udpAddr, _ = net.ResolveUDPAddr("udp", u.Config.Listen+":"+strconv.Itoa(u.Config.Port))
	fmt.Println(u.Config.Listen + ":" + strconv.Itoa(u.Config.Port))
	conn, err := net.ListenUDP("udp", u.udpAddr)
	if err != nil {
		log.Fatal("connect error", err)
	}
	for {
		var data = make([]byte, 8096)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println("READ error", err)
			continue
		}
		if n <= 0 {
			continue
		}
		fmt.Printf("[%v]:", remoteAddr)
		fmt.Println(data[:n])
		go u.handleConnection(remoteAddr, data[:n])

	}
	fmt.Println("go here")
	defer conn.Close()

}

func (u *UdpServer) handleConnection(addr *net.UDPAddr, message []byte) {
	fmt.Println(addr.String() + " connected!")
	//webSocketTunnel,err:= u.tunnelPool.get()
	//if err != nil {
	//	log.Println("get tunnel error ",err)
	//	return
	//}
	//webSocketTunnel.CreateUdpStreamOrSend(addr ,message)
}
