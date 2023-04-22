package topology

import (
	"context"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
	"sync"
	"time"
)

const (
	msgBuffSize      = 300 // Crazy ass number
	maxCumulativeMsg = 4   // Max of messages to send all together
)

// Topology represents the graph
// of all nodes in the network
type Topology struct {
	m       sync.RWMutex
	node    *maelstrom.Node
	network []string
	msgBuff chan float64
}

// New creates a new empty topology for a given node
func New(node *maelstrom.Node) *Topology {
	t := &Topology{
		node:    node,
		network: make([]string, 0),
		msgBuff: make(chan float64, msgBuffSize),
	}

	go t.broadcastFlusher()

	return t
}

// SetNetwork updates the network topology
func (t *Topology) SetNetwork() {
	t.m.Lock()
	defer t.m.Unlock()

	this := t.node.ID()
	nodes := t.node.NodeIDs()
	net := make([]string, len(nodes)-1)

	i := 0
	for _, n := range nodes {
		if n == this {
			continue
		}

		net[i] = n
		i++
	}

	t.network = net
}

// sendBroadcast emits a broadcast message to all nodes
func (t *Topology) sendBroadcast(messages []float64) {
	body := map[string]any{
		"type":           "broadcast",
		"batch_messages": messages,
	}

	t.m.RLock()
	defer t.m.RUnlock()

	// Spreads a broadcast message over
	// all nodes connected to sender
	for _, dest := range t.network {
		sendDest := dest
		go func() {
			for {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				// Message is ignored for now
				_, err := t.node.SyncRPC(ctx, sendDest, body)
				cancel() // Good practice...
				if _, ok := err.(*maelstrom.RPCError); ok {
					log.Printf("error while broadcasting to %v: %v", sendDest, err)
				} else if err == context.DeadlineExceeded {
					continue // Keep trying
				}

				return
			}
		}()
	}
}

// broadcastFlusher tries to send at most maxCumulativeMsg
// messages to all nodes in a single broadcast message
func (t *Topology) broadcastFlusher() {
	for {
		miniBuff := make([]float64, maxCumulativeMsg)

		// Waits for the first msg
		miniBuff[0] = <-t.msgBuff

		stop := false
		for i := 1; i < maxCumulativeMsg && !stop; i++ {
			select {
			case msg := <-t.msgBuff:
				miniBuff[i] = msg
			// Waits a bit to see if there are more messages
			case <-time.After(500 * time.Millisecond):
				stop = true
			}
		}

		go t.sendBroadcast(miniBuff)
	}
}

// Broadcast alerts that there's a new message to be sent.
// Note: messages are cumulative, i.e. 4 sent at a time in
// a field called batch_messages
func (t *Topology) Broadcast(msg float64) {
	go func() { t.msgBuff <- msg }()
}
