package service

import (
	"fmt"
	"github.com/ZYallers/rpcx-framework/helper"
	"github.com/go-redis/redis"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"github.com/soheilhy/cmux"
	"os"
	"time"
)

type Discovery struct {
	BasePath       string
	UpdateInterval time.Duration
	Addr           []string
}

type RPCXService struct {
	Env                string
	Version            string
	VersionKey         string
	Name               string
	HostName           string
	SystemIP           string
	Addr               string
	LogDir             string
	ErrorRobotToken    string
	GracefulRobotToken string
	Etcd               *Discovery
	Server             *server.Server
	SessionClient      func() *redis.Client
}

func (s *RPCXService) Serve() {
	defer func() {
		if err := recover(); err != nil {
			helper.SimpleMessage(fmt.Sprintf("%v", err), true, "panic")
		}
	}()

	s.Server = server.NewServer(
		//server.WithReadTimeout(10*time.Second),
		//server.WithWriteTimeout(15*time.Second),
		server.WithTCPKeepAlivePeriod(time.Minute),
	)
	s.Server.DisableJSONRPC = false
	s.Server.DisableHTTPGateway = true

	plugin := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + s.Addr,
		EtcdServers:    s.Etcd.Addr,
		BasePath:       s.Etcd.BasePath,
		UpdateInterval: s.Etcd.UpdateInterval,
	}
	if err := plugin.Start(); err != nil {
		panic(fmt.Sprintf("register etcdv3 plugin error: %s", err))
	} else {
		s.Server.Plugins.Add(plugin)
	}

	if err := RegisterFuncName(s, GetServices()); err != nil {
		panic(fmt.Sprintf("register function name error: %s", err))
	}

	s.Server.RegisterOnRestart(func(server *server.Server) {
		msg := fmt.Sprintf("%s service(%d) is restarting...", s.Name, os.Getpid())
		helper.SimpleMessage(msg, true, "info")
	})

	s.Server.RegisterOnShutdown(func(server *server.Server) {
		msg := fmt.Sprintf("%s service(%d) is shutting down...", s.Name, os.Getpid())
		helper.SimpleMessage(msg, true, "info")
	})

	msg := fmt.Sprintf("%s service(%d) is ready to serve", s.Name, os.Getpid())
	helper.SimpleMessage(msg, true, "info")
	if err := s.Server.Serve("tcp", s.Addr); err != nil && err != server.ErrServerClosed && err != cmux.ErrServerClosed {
		panic(fmt.Sprintf("%s service serve error: %v", s.Name, err))
	}
}
