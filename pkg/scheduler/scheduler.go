package scheduler

import (
	"fmt"
	"github.com/allen13/con-job/pkg/distributed"
	"github.com/allen13/con-job/pkg/distributed/etcdstore"
	"log"
	"strings"
	"sync"
)

type Scheduler struct {
	kvStore distributed.KeyValueStore
}

func Build() (scheduler Scheduler, err error) {
	scheduler = Scheduler{}
	etcdStore, err := etcdstore.Build()
	if err != nil {
		return
	}

	scheduler.kvStore = &etcdStore

	return
}

func (s *Scheduler) Start() {
	var watchWaitGroup sync.WaitGroup
	err := s.kvStore.RunForLeader()
	if err != nil {
		log.Println(err)
		//retryTimeout := config.GetEtcdTimeout()
		//log.Printf("attempting re-election in %s\n", retryTimeout/time.Second)
		//time.Sleep(retryTimeout)
		return
	}

	//elected leader start watching for events
	watchWaitGroup.Add(1)
	go s.kvStore.Watch("/nodes", s.onNodeKeyChange, &watchWaitGroup)

	watchWaitGroup.Add(1)
	go s.kvStore.Watch("/specifications", s.onSpecificationKeyChange, &watchWaitGroup)
	watchWaitGroup.Wait()

}

func (s *Scheduler) onNodeKeyChange(event distributed.KeyValueEvent) {
	separatedKeyPath := strings.Split(event.Key, "/")
	baseKeyName := separatedKeyPath[len(separatedKeyPath)-1]
	fmt.Println(baseKeyName)
	switch event.Type {
	case distributed.DELETE:
		s.kvStore.Delete(baseKeyName)
	case distributed.PUT:
		s.kvStore.Put(baseKeyName, event.Value)
	}
}

func (s *Scheduler) onSpecificationKeyChange(event distributed.KeyValueEvent) {
	separatedKeyPath := strings.Split(event.Key, "/")
	baseKeyName := separatedKeyPath[len(separatedKeyPath)-1]
	switch event.Type {
	case distributed.DELETE:
		s.kvStore.Delete(baseKeyName)
	case distributed.PUT:
		s.kvStore.Put(baseKeyName, event.Value)
	}
}
