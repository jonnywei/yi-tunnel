package client

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"yi-tunnel/common"
)

type UdpLocal struct {
	config     *common.Config
	udpAddr    *net.UDPAddr
	conn       *net.UDPConn
	tunnelPool *TunnelPool
}

func NewUdpLocal(config *common.Config, tunnelPool *TunnelPool) *UdpLocal {

	tcpLocal := UdpLocal{config: config}
	tcpLocal.tunnelPool = tunnelPool
	return &tcpLocal
}

func (u *UdpLocal) Listen() {
	u.udpAddr, _ = net.ResolveUDPAddr("udp", u.config.Local_address+":"+strconv.Itoa(u.config.Local_port))
	fmt.Println("udp listen:" + u.config.Local_address + ":" + strconv.Itoa(u.config.Local_port))
	conn, err := net.ListenUDP("udp", u.udpAddr)
	u.conn = conn
	if err != nil {
		log.Println("udp listen error", err)
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

func (u *UdpLocal) handleConnection(addr *net.UDPAddr, message []byte) {
	fmt.Println(addr.String() + " connected!")
	webSocketTunnel, err := u.tunnelPool.Get()
	if err != nil {
		log.Println("get tunnel error ", err)
		return
	}
	webSocketTunnel.CreateUdpStreamOrSend(u, addr, message)
}

func (u *UdpLocal) WriteToLocal(addr *net.UDPAddr, message []byte) {
	u.conn.WriteToUDP(message, addr)
}
