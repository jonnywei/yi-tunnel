package main

import (
	"github.com/jonnywei/yi_tunnel/client"
	"time"
)

func main() {

	config := client.LoadConfigFile("./config.json")
	client.RunClient(config)
	time.Sleep(time.Second * 13)
	client.Close()
	time.Sleep(time.Second * 10)
	client.RunClient(config)
	time.Sleep(time.Second * 13)
	client.Close()

	time.Sleep(time.Second * 10)
	client.RunClient(config)
	time.Sleep(time.Second * 13)
	client.Close()

	select {}
}
