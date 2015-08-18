package gopubsub

import (
	"sync"
)

type consumer func(payload []byte) error

type godis struct {
	channels map[string][]chan []byte
	sync.Mutex
}

type PubSub interface {
	Publish(channel string, payload []byte)
	Subscribe(channel string, c consumer)
}

func New() PubSub {
	return &godis{}
}

func (g *godis) Publish(channel string, payload []byte) {
	if chs, ok := g.channels[channel]; !ok {
		return
	}
	for ch := range chs {
		select {
		case ch <- payload:
			//
		}
	}
}

func (g *godis) Subscribe(channel string, c consumer) {
	if myChan, ok := g.channels[channel]; !ok {
		g.Lock()
		g.channels[channel] = make([]chan []byte, 0, 5)
		g.Unlock()
	}

	g.Lock()
	myChannel := make(chan []byte)
	g.channels[channel] = append(g.channels[channel], myChannel)
	g.Unlock()

	for pl := range myChannel {
		c(pl)
	}
}
