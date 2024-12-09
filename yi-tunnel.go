package yi_tunnel

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"yi_tunnel/client"
	"yi_tunnel/common"
)

func main() {

	var s common.Config

	configFlag := flag.String("c", "./config.json", "config file")

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

	tunnelPool := client.NewTunnelPool(&s)
	go func() {
		tcpLocal := client.NewTcpLocal(&s, tunnelPool)
		tcpLocal.Listen()

	}()
	udpLocal := client.NewUdpLocal(&s, tunnelPool)
	udpLocal.Listen()
}
