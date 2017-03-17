package scheduler

import (
	"github.com/allen13/con-job/pkg/distributed"
	"sync"
	"testing"
)

type MockKeyValueStore struct {
}

func (m *MockKeyValueStore) Watch(keypath string, keyValueEventCallback distributed.KeyValueEventCallback, watchWaitGroup sync.WaitGroup) {

}

func TestOnNodeKeyChange(t *testing.T) {

}
