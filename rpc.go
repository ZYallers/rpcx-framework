package framework

import (
	"errors"
	"github.com/ZYallers/golib/utils/logger"
	errors2 "github.com/ZYallers/rpcx-framework/errors"
	"github.com/ZYallers/rpcx-framework/helper/sender"
	"github.com/ZYallers/rpcx-framework/types"
	"github.com/go-redis/redis"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"
)

var rpc *types.Rpc

func NewService(options ...types.RpcOption) *types.Rpc {
	systemIP := SystemIP()
	if systemIP == "" || systemIP == "unknown" {
		panic("system ip is unknown or empty")
	}

	hostname := ServiceHostname()
	if hostname == "" {
		panic("system hostname is empty")
	}

	serviceName := ServiceName()
	if serviceName == "" {
		panic("service name is empty")
	}

	serviceLogDir := ServiceLogDir()
	if serviceLogDir == "" {
		panic("service log dir is empty")
	}

	discovery := ServiceDiscovery()
	if discovery == nil || discovery.BasePath == "" || len(discovery.Addr) == 0 {
		panic(errors2.ErrServiceDiscoveryNotMeeting)
	}

	rpc = &types.Rpc{
		Env:                ServiceMode(),
		Name:               serviceName,
		HostName:           hostname,
		SystemIP:           systemIP,
		LogDir:             serviceLogDir,
		Addr:               ServiceAddr(),
		Version:            ServiceVersion(),
		VersionKey:         ServiceVersionKey(),
		ErrorRobotToken:    ServiceErrorRobotToken(),
		GracefulRobotToken: ServiceGracefulRobotToken(),
		SqlRobotToken:      ServiceSqlRobotToken(),
		Etcd:               discovery,
		Server:             server.NewServer(),
	}

	for _, option := range options {
		if err := option(rpc); err != nil {
			panic(err)
		}
	}

	return rpc
}

func GetRpc() *types.Rpc { return rpc }

func WithSender() types.RpcOption {
	return func(s *types.Rpc) error {
		message := &types.Message{
			ErrorToken:    s.ErrorRobotToken,
			GracefulToken: s.GracefulRobotToken,
			SqlToken:      s.SqlRobotToken,
			Mode:          s.Env,
			Name:          s.Name,
			Addr:          s.Addr,
			Hostname:      s.HostName,
			SystemIP:      s.SystemIP,
			PublicIP:      PublicIP(),
		}
		types.InitMessage(message)
		s.Sender = message
		sender.Register(s.Sender)
		return nil
	}
}

func WithLogger() types.RpcOption {
	return func(s *types.Rpc) error {
		if s.LogDir == "" {
			return errors.New("service log dir is empty")
		}
		logger.SetLoggerDir(s.LogDir)
		log.SetLogger(types.NewLogger(s.Name, s.Sender))
		return nil
	}
}

func WithSessionFunc(fn func() *redis.Client) types.RpcOption {
	return func(s *types.Rpc) error {
		s.SessionFunc = fn
		return nil
	}
}
