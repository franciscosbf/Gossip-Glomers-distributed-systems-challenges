package topology

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// Topology represents the graph
// of all nodes in the network
type Topology struct {
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
		}

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
				// This "send" message doesn't give a s***
				// about response from the destination node
				_ = t.node.Send(sendDest, body)
			}()
		}
	}()
}
