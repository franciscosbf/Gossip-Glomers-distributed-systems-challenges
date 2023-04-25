package main

import (
	"encoding/json"
	"log"

	"maelstrom-gocounter/counter"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(n)
	c := counter.New(kv, n)

	n.Handle("add", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		delta := body["delta"].(float64)
		if err := c.Add(int(delta)); err != nil {
			return err
		}

		body["type"] = "add_ok"

		return n.Reply(msg, maelstrom.MessageBody{ Type: "add_ok" })
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		value, err := c.Get()
		if err != nil {
			return err
		}
		body["value"] = value

		body["type"] = "read_ok"

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
