package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jonnywei/yi_tunnel/common"
	"github.com/jonnywei/yi_tunnel/server"
	"log"
	"os"
)

func main() {

	var s common.ServerConfig

	configFlag := flag.String("sc", "./server_config.json", "config file")

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

	server := server.WebSocketServer{Config: s}

	server.Listen()

}
