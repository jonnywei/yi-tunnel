package client

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jonnywei/yi_tunnel/common"
	"log"
	"os"
)

type YiTunnelClient struct {
	tcpLocal *TcpLocal
	udpLocal *UdpLocal
	config   *common.Config
	filePath string
}

func NewYiTunnelClient(configFile string) *YiTunnelClient {
	ytc := &YiTunnelClient{filePath: configFile}
	return ytc
}

func (ytc *YiTunnelClient) LoadConfigFile() {
	var s common.Config

	configFlag := flag.String("c", ytc.filePath, "config file")

	flag.Parse()

	file, err := os.Open(*configFlag)
	if err != nil {
		log.Fatal("can't open config file", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&s)

	if err != nil {
		log.Fatal("can't decode config JSON: ", err)
	}
	fmt.Println(s)

	ytc.config = &s

}
func (ytc *YiTunnelClient) ListenAndServe() {
	tunnelPool := NewTunnelPool(ytc.config)
	var tcpLocal = NewTcpLocal(ytc.config, tunnelPool)
	var udpLocal = NewUdpLocal(ytc.config, tunnelPool)
	ytc.tcpLocal = tcpLocal
	ytc.udpLocal = udpLocal
	go tcpLocal.Listen()
	go udpLocal.Listen()
}

func (ytc *YiTunnelClient) Close() {
	ytc.tcpLocal.Close()
	ytc.udpLocal.Close()
}
