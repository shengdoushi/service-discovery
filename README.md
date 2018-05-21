# 用etcd来实现服务发现

本例为用etcd实现服务发现的简单实现。基于 etcd v3 实现.

# 安装

```
go get -u github.com/shengdoushi/service-discovery
```

# 服务注册:

流程:

1. 创建一个服务对象 Service
2. 调用 Start() 接口来注册
3. 调用 Stop() 来反注册

示例:

```golang
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("connect etcd fail:", err.Error())
		return
	}
	defer etcdClient.Close()

	// 1. 创建一个服务对象 
	s := NewService(etcdClient, "test/node/node1", "localhost:12312", &ServiceConfig{
		HeartbeatSeconds:4,
	})
	
	// 2. 注册服务， 返回一个关闭chnnel，在关闭后会放值
	stoppedCh, err := s.Start()
	if err != nil {
		t.Errorf("Service Start fail: %s", err.Error())
		return
	}
	
	// 3. 移除服务
	s.Stop()
```
