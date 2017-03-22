package distributed

import "sync"

type KeyValueStore interface {
	RunForLeader() (err error)
	//Watch a keypath and use a callback to respond to individual events.
	Watch(keypath string, keyValueEventCallback KeyValueEventCallback, watchWaitGroup *sync.WaitGroup)
	//Put a key into the key value store
	Put(key string, value string) (err error)
	//Put a key into the key value store but drop the key once the client disconnects or times out
	SoftPut(key string, value string) (err error)
	//Delete a key from the key value store
	Delete(key string) (err error)
}

type KeyValueEventType int32

const (
	PUT    KeyValueEventType = 0
	DELETE KeyValueEventType = 1
)

type KeyValueEvent struct {
	//Event type: either put or delete
	Type KeyValueEventType
	//Event key
	Key string
	//Event value
	Value string
}

type KeyValueEventCallback func(KeyValueEvent)
