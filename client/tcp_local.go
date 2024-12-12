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

type TcpLocal struct {
	config     *common.Config
	tcpAddr    *net.TCPAddr
	tunnelPool *TunnelPool
	quitChan   chan bool
	exitedChan chan bool
}

func NewTcpLocal(config *common.Config, tunnelPool *TunnelPool) *TcpLocal {

	tcpLocal := TcpLocal{config: config,
		tunnelPool: tunnelPool,
		quitChan:   make(chan bool),
		exitedChan: make(chan bool),
	}
	return &tcpLocal
}

func (tcpLocal *TcpLocal) Listen() error {

	tcpLocal.tcpAddr, _ = net.ResolveTCPAddr("tcp", tcpLocal.config.Local_address+":"+strconv.Itoa(tcpLocal.config.Local_port))
	conn, err := net.ListenTCP("tcp", tcpLocal.tcpAddr)
	if err != nil {
		log.Fatal("connect error", err)
		return err
	}
	fmt.Println("TCP Listen:" + tcpLocal.config.Local_address + ":" + strconv.Itoa(tcpLocal.config.Local_port))
	for {
		select {
		case <-tcpLocal.quitChan:
			log.Println("local server quiting...")
			conn.Close()
			close(tcpLocal.exitedChan)
			return nil
		default:
			conn.SetDeadline(time.Now().Add(time.Second * 5))
			c, err := conn.Accept()
			if err != nil {
				var netErr *net.OpError
				ok := errors.As(err, &netErr)
				if ok && netErr.Timeout() {
					log.Println("net timeout")
					continue
				}
				log.Fatal("accept error\n", err)
				return err
			}
			go tcpLocal.handleConnection(c)
		}
	}
	fmt.Println("tcp routine go here")
	defer conn.Close()

	return nil
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

// close tcplocal
func (tcpLocal *TcpLocal) Close() {
	tcpLocal.quitChan <- true

	<-tcpLocal.exitedChan
	log.Println("local server closed")
}
