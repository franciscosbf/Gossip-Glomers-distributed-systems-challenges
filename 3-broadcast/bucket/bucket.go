package bucket

import "sync"

// Bucket represents a collection of messages
type Bucket struct {
	m        sync.RWMutex
	messages map[float64]struct{}
}

// New creates a new collection
// of unique messages
func New() *Bucket {
	return &Bucket{
		messages: make(map[float64]struct{}),
	}
}

// Insert adds a new message
func (b *Bucket) Insert(msg float64) {
	b.m.Lock()
	defer b.m.Unlock()

	b.messages[msg] = struct{}{}
}

// InsertMany adds a collections of messages
func (b *Bucket) InsertMany(msg []float64) {
	b.m.Lock()
	defer b.m.Unlock()

	for _, m := range msg {
		b.messages[m] = struct{}{}
	}
}

// List returns a slice containing all messages
func (b *Bucket) List() []float64 {
	b.m.RLock()
	defer b.m.RUnlock()

	// Creates a slice of messages
	msgCpy := make([]float64, len(b.messages))
	i := 0
	for v := range b.messages {
		msgCpy[i] = v
		i++
	}

	return msgCpy
}
