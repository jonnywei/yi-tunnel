package client

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jonnywei/yi_tunnel/common"
	"log"
	"os"
)

var tcpLocal *TcpLocal
var udpLocal *UdpLocal

func LoadConfigFile(configFile string) common.Config {
	var s common.Config

	configFlag := flag.String("c", configFile, "config file")

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

	return s

}
func RunClient(s common.Config) {
	tunnelPool := NewTunnelPool(&s)
	tcpLocal = NewTcpLocal(&s, tunnelPool)
	udpLocal = NewUdpLocal(&s, tunnelPool)
	go tcpLocal.Listen()
	go udpLocal.Listen()
}

func Close() {
	tcpLocal.Close()
	udpLocal.Close()
}
