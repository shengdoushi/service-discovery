package discovery

import (
	"github.com/coreos/etcd/clientv3"
	"fmt"
	"time"
	"log"
)

func Example_basic(){
	// Output:
	// sleep..
	// stop..
	// stopped
	
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("connect etcd fail :", err.Error())
		return
	}
	defer etcdClient.Close()

	s := NewService(etcdClient, "test/node/node1", "localhost:12312", &ServiceConfig{
		HeartbeatSeconds:2,
	})
	stoppedCh, err := s.Start()
	if err != nil {
		log.Fatal("Service Start fail: ", err.Error())
		return
	}

	fmt.Println("sleep..")
	time.Sleep(5*time.Second)
	fmt.Println("stop..")
	s.Stop()
	select {
	case  <- stoppedCh:
		fmt.Println("stopped")
	case <- time.After(1*time.Second):
		log.Fatal("Service Stop fail")
	}
}
