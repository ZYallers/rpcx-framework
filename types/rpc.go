package types

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/smallnest/rpcx/server"
	"github.com/soheilhy/cmux"
)

type Discovery struct {
	UpdateInterval time.Duration
	BasePath       string
	Addr           []string
}

type Rpc struct {
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
	SqlRobotToken      string
	Etcd               *Discovery
	Server             *server.Server
	SessionFunc        func() *redis.Client
	Sender
}

type RpcOption func(s *Rpc) error

type ServerOption func(s *server.Server) error

func (s *Rpc) Serve(options ...ServerOption) {
	defer func() {
		if err := recover(); err != nil && s.Sender != nil {
			s.Sender.Graceful(err, true, "panic")
		}
	}()

	for _, option := range options {
		if err := option(s.Server); err != nil {
			panic(err)
		}
	}

	if s.Sender != nil {
		s.Server.RegisterOnRestart(func(server *server.Server) {
			msg := fmt.Sprintf("%s service(%d) is restarting", s.Name, os.Getpid())
			s.Sender.Graceful(msg, true, "info")
		})
		s.Server.RegisterOnShutdown(func(server *server.Server) {
			msg := fmt.Sprintf("%s service(%d) is shutting down", s.Name, os.Getpid())
			s.Sender.Graceful(msg, true, "info")
		})
		s.Sender.Graceful(fmt.Sprintf("%s service(%d) is ready to serve", s.Name, os.Getpid()), true, "info")
	}

	if err := s.Server.Serve("tcp", s.Addr); err != nil && err != server.ErrServerClosed && err != cmux.ErrServerClosed {
		panic(fmt.Sprintf("%s service serve error: %v", s.Name, err))
	}
}
