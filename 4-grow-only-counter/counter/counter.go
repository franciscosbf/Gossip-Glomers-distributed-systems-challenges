package counter

import (
	"context"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Inc struct {
	m sync.RWMutex
	kv *maelstrom.KV
	node *maelstrom.Node
}

// Create initializes the counter at zero.
// This doesn't givre a s*** about retries
func (i *Inc) Create() error {
	i.m.Lock()
	defer i.m.Unlock()

	ctx := context.Background()
	return i.kv.CompareAndSwap(ctx, i.node.ID(), 0, 0, true)
}

// Add increments the current value
func (i *Inc) Add(delta int) error {
	i.m.Lock()
	defer i.m.Unlock()

	v, err := i.kv.ReadInt(context.Background(), i.node.ID())
	if err != nil {
		return err
	}

	return i.kv.Write(context.Background(), i.node.ID(), v + delta)
}

// Get returns the current value
func (i *Inc) Get() (int, error) {
	i.m.RLock()
	defer i.m.RUnlock()

	// Simply sums all counters
	sum := 0
	for _, k := range i.node.NodeIDs() {
		ctx := context.Background()
		value, err := i.kv.ReadInt(ctx, k)
		if err != nil {
			log.Printf(
				"unexpected error while reading key %v at KV: %v",
				k, err,
			)
			continue	
		}
		sum += value
	}

	return sum, nil
}

// New creates a new Grow-only counter
func New(kv *maelstrom.KV, node *maelstrom.Node) *Inc {
	return &Inc{
		kv : kv,
		node: node,
	}
}

