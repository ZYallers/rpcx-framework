package service

import (
	"context"
	"github.com/ZYallers/rpcx-framework/define"
	"github.com/ZYallers/rpcx-framework/env"
	"github.com/syyongx/php2go"
	"reflect"
)

const stateActive = "state=active"

func RegisterFuncName(rs *RPCXService, services define.Restful) error {
	if err := registerHealthFunc(rs); err != nil {
		return err
	}
	if len(services) > 0 {
		if err := registerServiceMethod(rs, &services); err != nil {
			return err
		}
	}
	return nil
}

func registerHealthFunc(rs *RPCXService) error {
	return rs.Server.RegisterFunctionName(rs.Name, "health", func(ctx context.Context,
		args map[string]interface{}, reply *interface{}) error {
		*reply = "ok"
		return nil
	}, stateActive)
}

func registerServiceMethod(rs *RPCXService, services *define.Restful) error {
	for path, handlers := range *services {
		if err := rs.Server.RegisterFunctionName(rs.Name, path, dispatchHandler(rs, handlers), stateActive); err != nil {
			return err
		}
	}
	return nil
}

func dispatchHandler(rs *RPCXService, handlers []define.RestHandler) func(ctx context.Context, args map[string]interface{}, reply *interface{}) error {
	return func(ctx context.Context, args map[string]interface{}, reply *interface{}) error {
		argsVersion := rs.Version
		if ver, ok := args[rs.VersionKey].(string); ok && ver != "" {
			argsVersion = ver
		}
		if handler := versionCompare(&handlers, argsVersion); handler == nil {
			return env.ErrVersionCompare
		} else {
			v := reflect.ValueOf(handler.Service)
			ptr := reflect.New(v.Type().Elem())
			ptr.Elem().Set(v.Elem())
			sv := ptr.Interface().(define.IService)
			sv.Construct(rs, ctx, args, reply)
			if handler.Signed && !sv.SignCheck() {
				return env.ErrSignature
			}
			if handler.Logged && !sv.LoginCheck() {
				return env.ErrNeedLogin
			}
			result := ptr.MethodByName(handler.Method).Call(nil)
			if result[0].IsNil() {
				return nil
			}
			return result[0].Interface().(error)
		}
	}
}

func versionCompare(handlers *[]define.RestHandler, version string) *define.RestHandler {
	for _, handler := range *handlers {
		if handler.Version == "" || handler.Version == version {
			return &handler
		}
		if le := len(handler.Version); handler.Version[le-1:] == "+" {
			vs := handler.Version[0 : le-1]
			if version == vs {
				return &handler
			}
			if php2go.VersionCompare(version, vs, ">") {
				return &handler
			}
		}
	}
	return nil
}
