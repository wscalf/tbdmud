package net

import (
	"sync"

	"github.com/wscalf/tbdmud/internal/game"
)

type AggregateClientListener struct {
	listeners []game.ClientListener
}

func NewAggregateClientListener() *AggregateClientListener {
	return &AggregateClientListener{
		listeners: []game.ClientListener{},
	}
}

func (l *AggregateClientListener) AddListener(listener game.ClientListener) {
	l.listeners = append(l.listeners, listener)
}

func (l *AggregateClientListener) Listen() (chan game.Client, error) {
	listenerCount := len(l.listeners)
	clientChannels := make([]chan game.Client, 0, listenerCount)
	for _, listener := range l.listeners {
		clients, err := listener.Listen()
		if err != nil {
			return nil, err
		}

		clientChannels = append(clientChannels, clients)
	}

	fannedInClients := make(chan game.Client, 5)

	var wg sync.WaitGroup
	wg.Add(listenerCount)

	for _, clients := range clientChannels {
		go func(ch chan game.Client) {
			defer wg.Done()
			for client := range ch {
				fannedInClients <- client
			}
		}(clients)
	}

	go func() {
		wg.Wait()
		close(fannedInClients)
	}()

	return fannedInClients, nil
}

func (l *AggregateClientListener) LastError() error {
	for _, listener := range l.listeners {
		if err := listener.LastError(); err != nil {
			return err
		}
	}

	return nil
}
