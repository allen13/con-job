package scheduler

import (
	"github.com/allen13/con-job/pkg/config"
	"github.com/allen13/con-job/pkg/distributed"
	"github.com/allen13/con-job/pkg/distributed/etcdstore"
	"log"
	"sync"
	"time"
)

type ConJobScheduler struct {
	kvStore distributed.KeyValueStore
}

func Build() (scheduler ConJobScheduler, err error) {
	scheduler = ConJobScheduler{}
	etcdStore, err := etcdstore.Build()
	if err != nil {
		return
	}

	scheduler.kvStore = &etcdStore

	return
}

func (s *ConJobScheduler) Start() {
	var watchWaitGroup sync.WaitGroup
	for {
		err := s.kvStore.RunForLeader()
		if err != nil {
			log.Println(err)
			retryTimeout := config.GetEtcdTimeout()
			log.Printf("attempting re-election in %s\n", retryTimeout/time.Second)
			time.Sleep(retryTimeout)
			continue
		}

		//elected leader start watching for events
		watchWaitGroup.Add(1)
		go s.kvStore.Watch("/nodes", s.onNodeKeyChange, watchWaitGroup)

		watchWaitGroup.Add(1)
		go s.kvStore.Watch("/specifications", s.onSpecificationKeyChange, watchWaitGroup)

		watchWaitGroup.Wait()
	}

}

func (s *ConJobScheduler) onNodeKeyChange(event distributed.KeyValueEvent) {
	switch event.Type {
	case distributed.DELETE:
		//what to do when a node gets deleted
	case distributed.PUT:
		//what to do when a node gets added
	}
}

func (s *ConJobScheduler) onSpecificationKeyChange(event distributed.KeyValueEvent) {
	switch event.Type {
	case distributed.DELETE:
	//what to do when a specification gets deleted
	case distributed.PUT:
		//what to do when a specification gets added
	}
}
