package main

import (
	"github.com/jonnywei/yi_tunnel/client"
	"time"
)

func main() {
	ytc := client.NewYiTunnelClient("./config.json")
	ytc.LoadConfigFile()
	ytc.ListenAndServe()
	time.Sleep(time.Second * 13)
	ytc.Close()
	time.Sleep(time.Second * 10)
	ytc.ListenAndServe()
	time.Sleep(time.Second * 13)
	ytc.Close()

	time.Sleep(time.Second * 10)
	ytc.ListenAndServe()
	time.Sleep(time.Second * 13)
	ytc.Close()

	select {}
}
