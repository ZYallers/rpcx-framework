package client

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	framework "github.com/ZYallers/rpcx-framework"
	"github.com/ZYallers/rpcx-framework/consts"
	errors2 "github.com/ZYallers/rpcx-framework/errors"
	"github.com/ZYallers/rpcx-framework/helper/safe"
	"github.com/ZYallers/rpcx-framework/helper/sender"
	client2 "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
)

const (
	xClientDefaultTimeout = 15 * time.Second
)

var (
	discoveryDict = safe.NewDict()
	xClientDict   = safe.NewDict()
	failMode      = client.Failover
	selectMode    = client.RoundRobin
	xClientOption = client.Option{
		Retries:            3,                     // sets retries to send
		RPCPath:            share.DefaultRPCPath,  // sets for http connection
		ConnectTimeout:     time.Second,           // sets timeout for dialing
		SerializeType:      protocol.MsgPack,      // sets serialization type of payload
		CompressType:       protocol.None,         // sets decompression type
		BackupLatency:      10 * time.Millisecond, // is used for Failbackup mode, rpcx will sends another request if the first response doesn't return in BackupLatency time
		TCPKeepAlivePeriod: time.Minute,           // if it is zero we don't set keepalive
		IdleTimeout:        xClientDefaultTimeout, // ReadTimeout sets max idle time for underlying net.Conns
		GenBreaker: func() client.Breaker {
			// if failed 10 times, return error immediately, and will try to connect after 60 seconds
			return client.NewConsecCircuitBreaker(10, time.Minute)
		}, // is used to config CircuitBreaker
	}
)

func init() {
	switch framework.ServiceMode() {
	case consts.DevelopMode:
		failMode = client.Failfast
		selectMode = client.RandomSelect
	}
}

func XClient(service, serviceMethod string, args map[string]interface{}) (reply interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("xclient recover: %v", r)
		}
	}()

	sd := framework.ServiceDiscovery()
	if sd == nil || sd.BasePath == "" || len(sd.Addr) == 0 {
		err = errors2.ErrServiceDiscoveryNotMeeting
		return
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(10000))
	renew := n.Int64()%111 == 0
	xClient, err := getXClient(sd.BasePath, service, sd.Addr, renew)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), xClientDefaultTimeout)
	defer cancel()

	err = xClient.Call(ctx, serviceMethod, args, &reply)
	return
}

func getEtcdV3Discovery(basePath, service string, addr []string) (client.ServiceDiscovery, error) {
	if v, ok := discoveryDict.Get(service); ok {
		// log.Printf("loaded EtcdV3Discovery: %p, service: %s\n", v, service)
		return v.(client.ServiceDiscovery), nil
	}

	if v, loaded := discoveryDict.GetOrPutFunc(service, func(key string) (interface{}, error) {
		if dis, err := client2.NewEtcdV3Discovery(basePath, key, addr, false, nil); err != nil {
			return nil, fmt.Errorf("new etcd discovery error: %v", err)
		} else {
			return dis, nil
		}
	}); loaded {
		// log.Printf("put EtcdV3Discovery haved old value: %p, service: %s\n", v, service)
		return v.(client.ServiceDiscovery), nil
	} else {
		// log.Printf("new EtcdV3Discovery: %p, service: %s\n", v, service)
		return v.(client.ServiceDiscovery), nil
	}
}

func getXClient(basePath, service string, addr []string, renew bool) (client.XClient, error) {
	if v, ok := xClientDict.Get(service); ok {
		// log.Printf("loaded xclient: %p, service: %s\n", v, service)
		if renew {
			if ov, loaded := xClientDict.Delete(service); loaded {
				go func(s string, v interface{}) {
					defer safe.Defer()
					<-time.After(xClientDefaultTimeout)
					err := v.(client.XClient).Close()
					sender.Graceful(fmt.Sprintf("renew %s xclient: %v", s, err), true)
				}(service, ov)
			}
			if ov, loaded := discoveryDict.Delete(service); loaded {
				go func(s string, v interface{}) {
					defer safe.Defer()
					<-time.After(xClientDefaultTimeout + time.Second)
					v.(client.ServiceDiscovery).Close()
					sender.Graceful(fmt.Sprintf("renew %s discovery", s), true)
				}(service, ov)
			}
		} else {
			return v.(client.XClient), nil
		}
	}

	if v, loaded := xClientDict.GetOrPutFunc(service, func(key string) (interface{}, error) {
		if dis, err := getEtcdV3Discovery(basePath, key, addr); err != nil {
			return nil, err
		} else {
			return client.NewXClient(key, failMode, selectMode, dis, xClientOption), nil
		}
	}); loaded {
		// log.Printf("put xclient have old value: %p, service: %s\n", v, service)
		return v.(client.XClient), nil
	} else {
		// log.Printf("new xclient: %p, service: %s\n", v, service)
		return v.(client.XClient), nil
	}
}
