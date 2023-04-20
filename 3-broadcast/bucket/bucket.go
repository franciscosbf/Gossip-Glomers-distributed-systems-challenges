package bucket

import "sync"

// Bucket represents a collection of messages
type Bucket struct {
	m        sync.Mutex
	messages map[float64]struct{}
}

// New creates a new collection
// of unique messages
func New() *Bucket {
	return &Bucket{
		messages: make(map[float64]struct{}),
	}
}

// Insert adds a new message if it doesn't exist
func (b *Bucket) Insert(msg float64) {
	b.m.Lock()
	defer b.m.Unlock()

	if _, ok := b.messages[msg]; !ok {
		b.messages[msg] = struct{}{}
	}
}

// List returns a slice containing all messages
func (b *Bucket) List() []float64 {
	b.m.Lock()
	defer b.m.Unlock()

	// Creates a slice of messages
	msgCpy := make([]float64, len(b.messages))
	i := 0
	for v := range b.messages {
		msgCpy[i] = v
		i++
	}

	return msgCpy
}
