package restful

import (
	"context"
	"reflect"

	"github.com/ZYallers/rpcx-framework/errors"
	"github.com/ZYallers/rpcx-framework/types"
	"github.com/syyongx/php2go"
)

const stateActive = "state=active"

func RegisterFuncName(rs *types.Rpc, services types.Restful) error {
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

func registerHealthFunc(rs *types.Rpc) error {
	return rs.Server.RegisterFunctionName(rs.Name, "health", func(ctx context.Context,
		args map[string]interface{}, reply *interface{}) error {
		*reply = "ok"
		return nil
	}, stateActive)
}

func registerServiceMethod(rs *types.Rpc, services *types.Restful) error {
	for path, handlers := range *services {
		if err := rs.Server.RegisterFunctionName(rs.Name, path, dispatchHandler(rs, handlers), stateActive); err != nil {
			return err
		}
	}
	return nil
}

func dispatchHandler(rs *types.Rpc, handlers []types.RestHandler) func(ctx context.Context, args map[string]interface{}, reply *interface{}) error {
	return func(ctx context.Context, args map[string]interface{}, reply *interface{}) error {
		argsVersion := rs.Version
		if ver, ok := args[rs.VersionKey].(string); ok && ver != "" {
			argsVersion = ver
		}
		if handler := versionCompare(&handlers, argsVersion); handler == nil {
			return errors.ErrVersionCompare
		} else {
			v := reflect.ValueOf(handler.Service)
			ptr := reflect.New(v.Type().Elem())
			ptr.Elem().Set(v.Elem())
			sv := ptr.Interface().(types.IService)
			sv.Construct(rs, ctx, args, reply)
			if handler.Signed && !sv.SignCheck() {
				return errors.ErrSignature
			}
			if handler.Logged && !sv.LoginCheck() {
				return errors.ErrNeedLogin
			}
			result := ptr.MethodByName(handler.Method).Call(nil)
			if result[0].IsNil() {
				return nil
			}
			return result[0].Interface().(error)
		}
	}
}

func versionCompare(handlers *[]types.RestHandler, version string) *types.RestHandler {
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
