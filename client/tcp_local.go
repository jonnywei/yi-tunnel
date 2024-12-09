package client

import (
	"fmt"
	"github.com/jonnywei/yi_tunnel/common"
	"log"
	"net"
	"strconv"
)

type TcpLocal struct {
	config     *common.Config
	tcpAddr    *net.TCPAddr
	tunnelPool *TunnelPool
}

func NewTcpLocal(config *common.Config, tunnelPool *TunnelPool) *TcpLocal {

	tcpLocal := TcpLocal{config: config}
	tcpLocal.tunnelPool = tunnelPool
	return &tcpLocal
}

func (tcpLocal *TcpLocal) Listen() {

	tcpLocal.tcpAddr, _ = net.ResolveTCPAddr("tcp", tcpLocal.config.Local_address+":"+strconv.Itoa(tcpLocal.config.Local_port))

	fmt.Println(tcpLocal.config.Local_address + ":" + strconv.Itoa(tcpLocal.config.Local_port))
	conn, err := net.ListenTCP("tcp", tcpLocal.tcpAddr)
	if err != nil {
		log.Fatal("connect error", err)
	}

	for {
		c, err := conn.Accept()
		if err != nil {
			log.Fatal("accept error", err)
			break
		}
		go tcpLocal.handleConnection(c)

	}
	fmt.Println("go here")
	defer conn.Close()

}

func (tcpLocal *TcpLocal) handleConnection(con net.Conn) {
	fmt.Println(con.RemoteAddr().String() + " connected!")
	webSocketTunnel, err := tcpLocal.tunnelPool.Get()
	if err != nil {
		log.Println("get tunnel error ", err)
		con.Close()
		return
	}
	webSocketTunnel.CreateStream(con)
}
