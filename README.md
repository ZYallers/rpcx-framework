# RPCX Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/ZYallers/rpcx-framework)](https://goreportcard.com/report/github.com/ZYallers/rpcx-framework)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/ZYallers/rpcx-framework.svg?branch=master)](https://travis-ci.org/ZYallers/rpcx-framework) 
[![Foundation](https://img.shields.io/badge/Golang-Foundation-green.svg)](http://golangfoundation.org) 
[![GoDoc](https://pkg.go.dev/badge/github.com/ZYallers/rpcx-framework?status.svg)](https://pkg.go.dev/github.com/ZYallers/rpcx-framework?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/ZYallers/rpcx-framework/-/badge.svg)](https://sourcegraph.com/github.com/ZYallers/rpcx-framework?badge)
[![Release](https://img.shields.io/github/release/ZYallers/rpcx-framework.svg?style=flat-square)](https://github.com/ZYallers/rpcx-framework/releases)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/ZYallers/rpcx-framework)](https://www.tickgit.com/browse?repo=github.com/ZYallers/rpcx-framework)
[![goproxy.cn](https://goproxy.cn/stats/github.com/ZYallers/rpcx-framework/badges/download-count.svg)](https://goproxy.cn)

> An RPC microservices framework based on rpcx. 
>
> Features: simple and easy to use, ultra fast and efficient, powerful, service discovery, service governance, service layering, version control, routing label registration.
>
> Best microservices framework in Go, like alibaba Dubbo, but with more features, Scale easily. Try it. Test it. If you feel it's better, use it! 
>
> Java有Dubbo, Golang有RPCX!

# Installation
To install rpcx-framework package, you need to install Go and set your Go workspace first.

1. The first need Go installed (version 1.11+ is required), then you can use the below Go command to install rpcx-framework.
```bash
$ go get -u github.com/ZYallers/rpcx-framework
```

2. Import it in your code:
```go 
import "github.com/ZYallers/rpcx-framework" 
```

# Examples

The below is a simple example.

```go
package main

import (
	"github.com/smallnest/rpcx/log"
	framework "gitlab.sys.hxsapp.net/hxs/rpcx-framework"
)

func init() {
	framework.LoadConfig()
}

func main() {
	//share.Trace = true
	s := framework.NewService()
	log.Infof("Service-> %+v; Etcd-> %+v", *s, *(s.Etcd))
	s.Serve()
}

```

# How to deploy and run?
Copy the boot script "script / bootstrap. Sh" to the root directory of your project, and then execute it; If successful, you will see the following information:
```
 current path: /Users/cloud/gopath_hxsapp/rpcx-demo 
 download produce.sh 
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  7791  100  7791    0     0    221      0  0:00:35  0:00:35 --:--:--  1904
 download produce.sh(/Users/cloud/gopath_hxsapp/rpcx-demo/./bin/produce.sh) finished 
 service config: 
 ServiceName: rpcx-demo 
 ServiceAddr: 172.18.28.123:8978 
 LogDir: /apps/logs/go/rpcx-demo 
 Operation: 
     status                                  View service status
     sync                                    Synchronization service vendor resources
     build                                   Compile and generate service program
     reload                                  Smooth restart service
     quit                                    Stop service
     help                                    View help information for the help command
 For more information about an action, use the help command to view it
```
At the same time, after successful execution, it will create a new `bin` directory in your current directory, 
and generate a service compilation and deployment script `produce sh`

Execute the deployment script `./bin/produce.sh help`, it will tell you what to do next.
```
 service config: 
 ServiceName: rpcx-demo 
 ServiceAddr: 172.18.28.123:8978 
 LogDir: /apps/logs/go/rpcx-demo 
 Operation: 
     status                                  View service status
     sync                                    Synchronization service vendor resources
     build                                   Compile and generate service program
     reload                                  Smooth restart service
     quit                                    Stop service
     help                                    View help information for the help command
 For more information about an action, use the help command to view it
```

# Feature
An RPC service framework based on rpcx (fast, easy-to-use but powerful RPC service governance framework of go language).

- easy to use
- super fast and efficient
- powerful
- service discovery
- service governance
- service layering
- version control
- routing label registration.

### RPCX
rpcx is a RPC framework like [Alibaba Dubbo](http://dubbo.io/) and [Weibo Motan](https://github.com/weibocom/motan).

**rpcx** is created for targets:
1. **Simple**: easy to learn, easy to develop, easy to intergate and easy to deploy
2. **Performance**: high perforamnce (>= grpc-go)
3. **Cross-platform**: support _raw slice of bytes_, _JSON_, _Protobuf_ and _MessagePack_. Theoretically it can be used with java, php, python, c/c++, node.js, c# and other platforms
4. **Service discovery and service governance**: support zookeeper, etcd and consul.


It contains below features
- Support raw Go functions. There's no need to define proto files.
- Pluggable. Features can be extended such as service discovery, tracing.
- Support TCP, HTTP, [QUIC](https://en.wikipedia.org/wiki/QUIC) and [KCP](https://github.com/skywind3000/kcp)
- Support multiple codecs such as JSON, Protobuf, [MessagePack](https://msgpack.org/index.html) and raw bytes.
- Service discovery. Support peer2peer, configured peers, [zookeeper](https://zookeeper.apache.org), [etcd](https://github.com/coreos/etcd), [consul](https://www.consul.io) and [mDNS](https://en.wikipedia.org/wiki/Multicast_DNS).
- Fault tolerance：Failover, Failfast, Failtry.
- Load banlancing：support Random, RoundRobin, Consistent hashing, Weighted, network quality and Geography.
- Support Compression.
- Support passing metadata.
- Support Authorization.
- Support heartbeat and one-way request.
- Other features: metrics, log, timeout, alias, circuit breaker.
- Support bidirectional communication.
- Support access via HTTP so you can write clients in any programming languages.
- Support API gateway.
- Support backup request, forking and broadcast.

# Reference
- RPCX：https://rpcx.io
- Redis Command：http://redis.cn/commands.html
- GORM：https://gorm.io/zh_CN/docs
- Docker：http://www.dockerinfo.net/document

# License
Released under the [MIT License](https://github.com/ZYallers/rpcx-framework/blob/master/LICENSE)
