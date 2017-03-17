package main

import (
	"fmt"
	"log"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"golang.org/x/net/context"
	"sync"
	"time"
)

var endpoints = []string{"http://localhost:2379"}
var dialTimeout = time.Second * 5
var wg sync.WaitGroup

func main() {
	wg.Add(3)
	go electScheduler("s1")
	go electScheduler("s2")
	go electScheduler("s3")
	wg.Wait()
}

func electScheduler(id string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	//continually run for leader
	for {
		fmt.Printf("scheduler %s running for election\n", id)
		clientSession, err := concurrency.NewSession(cli, concurrency.WithTTL(5))
		if err != nil {
			log.Fatal(err)
		}
		election := concurrency.NewElection(clientSession, "/poller/s1/{asdf}")
		err = election.Campaign(context.Background(), id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s elected leader\n", id)

		timer := time.NewTimer(time.Second * 5)
		<-timer.C
		clientSession.Close()
		fmt.Printf("%s lost leadership\n", id)
	}
	wg.Done()
}
