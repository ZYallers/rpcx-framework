package framework

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/ZYallers/rpcx-framework/helper/restful"
	"github.com/ZYallers/rpcx-framework/types"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
)

// WithTLSConfig sets tls.Config.
func WithTLSConfig(cfg *tls.Config) types.ServerOption {
	return func(s *server.Server) error {
		server.WithTLSConfig(cfg)(s)
		return nil
	}
}

// WithReadTimeout sets readTimeout.
func WithReadTimeout(readTimeout time.Duration) types.ServerOption {
	return func(s *server.Server) error {
		server.WithReadTimeout(readTimeout)(s)
		return nil
	}
}

// WithWriteTimeout sets writeTimeout.
func WithWriteTimeout(writeTimeout time.Duration) types.ServerOption {
	return func(s *server.Server) error {
		server.WithWriteTimeout(writeTimeout)(s)
		return nil
	}
}

// WithTCPKeepAlivePeriod sets tcp keepalive period.
func WithTCPKeepAlivePeriod(period time.Duration) types.ServerOption {
	return func(s *server.Server) error {
		server.WithTCPKeepAlivePeriod(period)(s)
		return nil
	}
}

func WithDisableJSONRPC(disable bool) types.ServerOption {
	return func(s *server.Server) error {
		s.DisableJSONRPC = disable
		return nil
	}
}

func WithDisableHTTPGateway(disable bool) types.ServerOption {
	return func(s *server.Server) error {
		s.DisableHTTPGateway = disable
		return nil
	}
}

func WithEtcdV3Plugin(addr string, d *types.Discovery) types.ServerOption {
	return func(s *server.Server) error {
		plugin := &serverplugin.EtcdV3RegisterPlugin{
			ServiceAddress: "tcp@" + addr,
			EtcdServers:    d.Addr,
			BasePath:       d.BasePath,
			UpdateInterval: d.UpdateInterval,
		}
		if err := plugin.Start(); err != nil {
			return fmt.Errorf("etcdv3 plugin register error: %s", err)
		}
		s.Plugins.Add(plugin)
		return nil
	}
}

func WithFunction(rpc *types.Rpc, rest types.Restful) types.ServerOption {
	return func(s *server.Server) error {
		if err := restful.RegisterFuncName(rpc, rest); err != nil {
			return fmt.Errorf("register function error: %s", err)
		}
		return nil
	}
}

func WithOnRestart(f func(server *server.Server)) types.ServerOption {
	return func(s *server.Server) error {
		s.RegisterOnRestart(f)
		return nil
	}
}

func WithOnShutdown(f func(server *server.Server)) types.ServerOption {
	return func(s *server.Server) error {
		s.RegisterOnShutdown(f)
		return nil
	}
}
