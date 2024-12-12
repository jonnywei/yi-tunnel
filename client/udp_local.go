package client

import (
	"errors"
	"fmt"
	"github.com/jonnywei/yi_tunnel/common"
	"log"
	"net"
	"strconv"
	"time"
)

type UdpLocal struct {
	config     *common.Config
	udpAddr    *net.UDPAddr
	conn       *net.UDPConn
	tunnelPool *TunnelPool
	quitChan   chan bool
	exitedChan chan bool
}

func NewUdpLocal(config *common.Config, tunnelPool *TunnelPool) *UdpLocal {

	tcpLocal := UdpLocal{config: config,
		quitChan:   make(chan bool),
		exitedChan: make(chan bool),
	}
	tcpLocal.tunnelPool = tunnelPool
	return &tcpLocal
}

func (u *UdpLocal) Listen() {
	u.udpAddr, _ = net.ResolveUDPAddr("udp", u.config.Local_address+":"+strconv.Itoa(u.config.Local_port))
	fmt.Println("UDP Listen:" + u.config.Local_address + ":" + strconv.Itoa(u.config.Local_port))
	conn, err := net.ListenUDP("udp", u.udpAddr)
	u.conn = conn
	if err != nil {
		log.Println("udp listen error", err)
	}
	for {
		select {
		case <-u.quitChan:
			u.conn.Close()
			close(u.exitedChan)
			return
		default:
			var data = make([]byte, 8096)
			conn.SetDeadline(time.Now().Add(time.Second * 5))
			n, remoteAddr, err := conn.ReadFromUDP(data)
			if err != nil {
				log.Println("READ error", err)
				var netErr *net.OpError
				ok := errors.As(err, &netErr)
				if ok && netErr.Timeout() {
					log.Println("udp net timeout")
					continue
				}
				continue
			}
			if n <= 0 {
				continue
			}
			fmt.Printf("[%v]:", remoteAddr)
			fmt.Println(data[:n])
			go u.handleConnection(remoteAddr, data[:n])
		}
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

// close UdpLocal
func (u *UdpLocal) Close() {
	u.quitChan <- true
	<-u.exitedChan
	log.Println("udp local server closed")
}
