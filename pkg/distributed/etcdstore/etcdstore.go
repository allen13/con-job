package etcdstore

import (
	"context"
	"fmt"
	"github.com/allen13/con-job/pkg/config"
	"github.com/allen13/con-job/pkg/distributed"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"log"
	"sync"
	"time"
)

type EtcdStore struct {
	etcdClient *clientv3.Client
	leaseId    clientv3.LeaseID
}

func Build() (etcd EtcdStore, err error) {
	etcd = EtcdStore{}

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEtcdEndpoints(),
		DialTimeout: config.GetEtcdTimeout(),
	})
	if err != nil {
		return
	}
	etcd.etcdClient = etcdClient

	return
}

func (e *EtcdStore) RunForLeader() (err error) {
	hostname := config.GetHostname()

	fmt.Printf("scheduler %s running for election\n", hostname)

	clientSession, err := concurrency.NewSession(e.etcdClient, concurrency.WithTTL(5))
	if err != nil {
		return
	}

	election := concurrency.NewElection(clientSession, "scheduler")
	err = election.Campaign(context.Background(), hostname)
	if err != nil {
		return
	}

	log.Printf("%s elected leader\n", hostname)

	return
}

func (e *EtcdStore) Watch(keypath string, keyValueEventCallback distributed.KeyValueEventCallback, watchWaitGroup sync.WaitGroup) {
	watchChannel := e.etcdClient.Watch(context.Background(), keypath)
	for watchResponse := range watchChannel {
		for _, event := range watchResponse.Events {
			kvEvent := distributed.KeyValueEvent{
				Key:   string(event.Kv.Key),
				Value: string(event.Kv.Value),
			}
			switch event.Type {
			case mvccpb.DELETE:
				kvEvent.Type = distributed.DELETE

			case mvccpb.PUT:
				kvEvent.Type = distributed.PUT
			}

			keyValueEventCallback(kvEvent)
		}
	}
	watchWaitGroup.Done()
}

func (e *EtcdStore) Put(key string, value string) (err error) {
	_, err = e.etcdClient.Put(context.Background(), key, value)
	return
}

func (e *EtcdStore) SoftPut(key string, value string) (err error) {
	leaseId, err := e.getLeaseId()
	if err != nil {
		return
	}
	_, err = e.etcdClient.Put(context.Background(), key, value, clientv3.WithLease(leaseId))
	return
}

func (e *EtcdStore) getLeaseId() (clientv3.LeaseID, error) {
	if e.leaseId == 0 {
		return e.leaseId, nil
	}

	leaseGrantResponse, err := e.etcdClient.Grant(context.Background(), 10)
	if err != nil {
		return 0, err
	}

	e.leaseId = leaseGrantResponse.ID

	go e.keepLeaseAlive()

	return e.leaseId, nil
}

func (e *EtcdStore) keepLeaseAlive() {
	ticker := time.NewTicker(time.Second * 10)
	for _ := range ticker.C {
		e.etcdClient.KeepAlive(context.Background(), e.leaseId)
	}
}
