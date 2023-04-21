package topology

import (
	"context"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
	"sync"
	"time"
)

// Topology represents the graph
// of all nodes in the network
type Topology struct {
	m       sync.Mutex
	node    *maelstrom.Node
	network map[string][]string
}

// New creates a new empty topology for a given node
func New(node *maelstrom.Node) *Topology {
	return &Topology{
		node:    node,
		network: make(map[string][]string),
	}
}

// SetNetwork updates the network topology
func (t *Topology) SetNetwork(net map[string][]string) {
	t.m.Lock()
	defer t.m.Unlock()

	t.network = net
}

// Broadcast sends a given message to all neighbors.
// The requester node won't receive this message
func (t *Topology) Broadcast(requester string, msg any) {
	// Broadcasting is done by a background task since
	// we don't need to ensure this with the requester
	go func() {
		nodeId := t.node.ID()
		body := map[string]any{
			"type":    "broadcast",
			"message": msg,
			"_test":   1,
		}

		t.m.Lock()
		defer t.m.Unlock()

		// Spreads a broadcast message over
		// all nodes connected to sender
		for _, dest := range t.network[nodeId] {
			// The broadcast requester won't
			// receive what he asked for
			if dest == requester {
				continue
			}

			sendDest := dest
			go func() {
				for {
					ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
					// Message is ignored for now
					_, err := t.node.SyncRPC(ctx, sendDest, body)
					cancel() // Good practice...
					if _, ok := err.(*maelstrom.RPCError); ok {
						log.Printf("error while broadcasting to %v: %v", sendDest, err)
					} else if err == context.DeadlineExceeded {
						time.Sleep(1 * time.Second)

						continue // Keep trying
					}

					return
				}
			}()
		}
	}()
}
