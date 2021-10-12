package discovery

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"context"
)

type Service struct {
	nodeKey string
	nodeAddress string
	stopCh chan struct{}
	etcdClient *clientv3.Client
	config *ServiceConfig
}

type ServiceConfig struct {
	HeartbeatSeconds int32		// default 5 seconds
}

// Create a Service
func NewService(etcdClient* clientv3.Client, nodeKey string, nodeAddress string, config *ServiceConfig)(*Service){
	return &Service{
		nodeKey:nodeKey,
		nodeAddress:nodeAddress,
		stopCh:make(chan struct{}, 2),
		etcdClient:etcdClient,
		config:config,
	}
}

// Start Service: register to ETCD, and start heatbeat
func (service *Service)Start()(<- chan struct{}, error){
	heatbeatSeconds := int64(5)
	if service.config != nil && service.config.HeartbeatSeconds > 0 {
		heatbeatSeconds = int64(service.config.HeartbeatSeconds)
	}
	stoppedCh := make(chan struct{}, 1)
	
	// register
	grantResp, err := service.etcdClient.Grant(context.TODO(), heatbeatSeconds)
	if err != nil {
		return nil, err
	}
	grantId := grantResp.ID

	_, err = service.etcdClient.Put(context.TODO(), service.nodeKey, service.nodeAddress, clientv3.WithLease(grantId))
	if err != nil {
		return nil, err
	}

	keepAliveResp, err := service.etcdClient.KeepAlive(context.TODO(), grantId)
	if err != nil {
		return nil, err
	}

	// heartbeat
	go func(){
		defer func() {
			stoppedCh <- struct{}{}
		}()
	LOOP:
		for {
			select {
			case <- service.stopCh:
				service.etcdClient.Revoke(context.TODO(), grantId)
				break LOOP
			case <-keepAliveResp: // todo : if fail?
				break
			}
		}
	}()
	return stoppedCh, nil
}

// Stop a Service: delete from ETCD
func (service *Service)Stop(){
	service.stopCh <- struct{}{}
}
