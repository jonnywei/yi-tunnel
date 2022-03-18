package client

import (
	"fmt"
	"sync"
	"yi-tunnel/common"
)

type TunnelPool struct {
	tunnelMap map[string]*WebSocketTunnel
	config    *common.Config
	sync.Mutex
}

func NewTunnelPool(config *common.Config) *TunnelPool {
	return &TunnelPool{
		config:    config,
		tunnelMap: make(map[string]*WebSocketTunnel),
	}

}

func (p *TunnelPool) Get() (*WebSocketTunnel, error) {

	p.Lock()
	defer p.Unlock()
	for k, t := range p.tunnelMap {
		if t.IsClosed() {
			delete(p.tunnelMap, k)
		} else if t.StreamCount() < p.config.Stream_count_per_channel {
			return t, nil
		}
	}
	fmt.Println("new tunnel")
	tunnel := NewWebSocketTunnel(p.config)
	err := tunnel.Open()
	if err != nil {
		return nil, err
	}
	var name = tunnel.Name
	p.tunnelMap[name] = tunnel
	return tunnel, err
}
