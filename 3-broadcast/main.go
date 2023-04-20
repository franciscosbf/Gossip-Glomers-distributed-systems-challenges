package main

import (
	"encoding/json"
	"log"

	"maelstrom-broadcast/bucket"
	"maelstrom-broadcast/topology"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	msgBucket := bucket.New()
	network := topology.New(n)

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Extracts the message and stores it
		content := body["message"].(float64)
		msgBucket.Insert(content)

		// Broadcasts this msg to all connected nodes
		msgId := body["msg_id"]
		network.Broadcast(msg.Src, msgId, content)

		// Removes message entry from reply body
		delete(body, "message")

		// Sets the response confirmation
		body["type"] = "broadcast_ok"

		return n.Reply(msg, body)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Sets all stored messages
		body["messages"] = msgBucket.List()

		// Sets the response confirmation
		body["type"] = "read_ok"

		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Extracts the topology and sets it to the network
		rawNet := body["topology"].(map[string]any)
		newNet := make(map[string][]string)
		for k, v := range rawNet {
			rawNeighbors := v.([]any)
			neighbors := make([]string, len(rawNeighbors))
			for i := 0; i < len(rawNeighbors); i++ {
				neighbors[i] = rawNeighbors[i].(string)
			}
			newNet[k] = neighbors
		}
		network.SetNetwork(newNet)

		// Removes topology entry from reply body
		delete(body, "topology")

		// Sets the response confirmation
		body["type"] = "topology_ok"

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
