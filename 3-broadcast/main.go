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

		// Extracts the message(s) and stores it(them)
		if content, ok := body["message"]; ok {
			value := content.(float64)
			msgBucket.Insert(value)
			network.Broadcast(value)
		} else { // We assume that batch_messages exists
			content = body["batch_messages"]
			rawValues := content.([]any)
			values := make([]float64, len(rawValues))
			for i := 0; i < len(rawValues); i++ {
				values[i] = rawValues[i].(float64)
			}
			msgBucket.InsertMany(values)
			// Those aren't forwarded
		}

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

		network.SetNetwork()

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
