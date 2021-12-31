package client

import (
	"context"
	framework "github.com/ZYallers/rpcx-framework"
	"github.com/ZYallers/rpcx-framework/env"
	"github.com/ZYallers/zgin/libraries/tool"
	client2 "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
	"sync"
	"time"
)

var (
	xClientMap   sync.Map
	failMode     = client.Failover
	selectMode   = client.RoundRobin
	clientOption = client.Option{
		Retries:            3,
		RPCPath:            share.DefaultRPCPath,
		ConnectTimeout:     time.Second,
		SerializeType:      protocol.MsgPack,
		CompressType:       protocol.None,
		BackupLatency:      10 * time.Millisecond,
		TCPKeepAlivePeriod: time.Minute, // if it is zero we don't set keepalive
		IdleTimeout:        time.Minute, // ReadTimeout sets max idle time for underlying net.Conns
		GenBreaker: func() client.Breaker {
			// if failed 10 times, return error immediately, and will try to connect after 60 seconds
			return client.NewConsecCircuitBreaker(10, 60*time.Second)
		},
	}
)

func init() {
	switch framework.ServiceMode() {
	case env.DevelopMode:
		failMode = client.Failfast
		selectMode = client.RandomSelect
	}
}

//  XClient ...
//  @author Cloud|2021-12-23 15:14:02
//  @param serviceName string ...
//  @param serviceMethod string ...
//  @param args map[string]interface{} ...
//  @return interface{} ...
//  @return error ...
func XClient(serviceName, serviceMethod string, args map[string]interface{}) (interface{}, error) {
	share.Trace = false
	if val, ok := args["trace"]; ok && val.(string) == "on" {
		share.Trace = true
	}

	var xClient client.XClient
	if val, ok := xClientMap.Load(serviceName); ok {
		xClient = val.(client.XClient)
	} else {
		basePath := framework.DiscoveryBasePath()
		addr := framework.DiscoveryAddress()
		d, _ := client2.NewEtcdV3Discovery(basePath, serviceName, addr, false, nil)
		xClient = client.NewXClient(serviceName, failMode, selectMode, d, clientOption)
		xClientMap.Store(serviceName, xClient)
	}

	var reply interface{}
	ctx, cancel := context.WithTimeout(context.Background(), tool.DefaultHttpClientTimeout)
	defer cancel()
	return reply, xClient.Call(ctx, serviceMethod, args, &reply)
}
