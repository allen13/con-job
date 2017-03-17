package distributed

import "sync"

type KeyValueStore interface {
	RunForLeader() (err error)
	Watch(keypath string, keyValueEventCallback KeyValueEventCallback)
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
