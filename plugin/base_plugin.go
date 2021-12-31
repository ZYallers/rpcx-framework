package plugin

import (
	"context"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/protocol"
	"net"
)

type BasePlugin struct{}

func (p *BasePlugin) Register(name string, rcvr interface{}, metadata string) error {
	log.Infof("BasePlugin.Register-> name:%s, rcvr:%+v, metadata:%s", name, rcvr, metadata)
	return nil
}

func (p *BasePlugin) Unregister(name string) error {
	log.Infof("BasePlugin.Unregister-> name:%s", name)
	return nil
}

func (p *BasePlugin) RegisterFunction(serviceName, fname string, fn interface{}, metadata string) error {
	log.Infof("BasePlugin.RegisterFunction-> serviceName:%s, fname:%s, fn:%+v, metadata:%s",
		serviceName, fname, fn, metadata)
	return nil
}

func (p *BasePlugin) HandleConnAccept(conn net.Conn) (net.Conn, bool) {
	log.Infof("BasePlugin.HandleConnAccept-> conn:%+v, LocalAddr:%s, RemoteAddr:%s",
		conn, conn.LocalAddr().String(), conn.RemoteAddr().String())
	return conn, true
}

func (p *BasePlugin) HandleConnClose(conn net.Conn) bool {
	log.Infof("BasePlugin.HandleConnClose-> conn:%+v, LocalAddr:%s, RemoteAddr:%s",
		conn, conn.LocalAddr().String(), conn.RemoteAddr().String())
	return true
}

func (p *BasePlugin) PreReadRequest(ctx context.Context) error {
	log.Infof("BasePlugin.PreReadRequest-> ctx:%+v", ctx)
	return nil
}

func (p *BasePlugin) PostReadRequest(ctx context.Context, r *protocol.Message, e error) error {
	log.Infof("BasePlugin.PostReadRequest-> ctx:%+v, r:%+v, e:%+v", ctx, r, e)
	return nil
}

func (p *BasePlugin) PreHandleRequest(ctx context.Context, r *protocol.Message) error {
	log.Infof("BasePlugin.PreHandleRequest-> ctx:%+v, r:%+v", ctx, r)
	return nil
}

func (p *BasePlugin) PreCall(ctx context.Context, serviceName, methodName string, args interface{}) (interface{}, error) {
	log.Infof("BasePlugin.PreCall-> ctx:%+v, serviceName:%s, methodName:%s, args:%+v", ctx, serviceName, methodName, args)
	return args, nil
}

func (p *BasePlugin) PostCall(ctx context.Context, serviceName, methodName string, args, reply interface{}) (interface{}, error) {
	log.Infof("BasePlugin.PostCall-> ctx:%+v, serviceName:%s, methodName:%s, args:%+v, reply:%+v", ctx, serviceName, methodName, args, reply)
	return reply, nil
}

func (p *BasePlugin) PreWriteResponse(ctx context.Context, req *protocol.Message, resp *protocol.Message, err error) error {
	log.Infof("BasePlugin.PreWriteResponse-> ctx:%+v, req:%+v, resp:%+v, err:%+v", ctx, req, resp, err)
	return nil
}

func (p *BasePlugin) PostWriteResponse(ctx context.Context, req *protocol.Message, resp *protocol.Message, err error) error {
	log.Infof("BasePlugin.PostWriteResponse-> ctx:%+v, req:%+v, resp:%+v, err:%+v", ctx, req, resp, err)
	return nil
}

func (p *BasePlugin) PreWriteRequest(ctx context.Context) error {
	log.Infof("BasePlugin.PreWriteRequest-> ctx:%+v", ctx)
	return nil
}

func (p *BasePlugin) PostWriteRequest(ctx context.Context, r *protocol.Message, e error) error {
	log.Infof("BasePlugin.PostWriteRequest-> ctx:%+v, r:%+v, e:%+v", ctx, r, e)
	return nil
}

func (p *BasePlugin) HeartbeatRequest(ctx context.Context, req *protocol.Message) error {
	log.Infof("BasePlugin.HeartbeatRequest-> ctx:%+v, req:%+v", ctx, req)
	return nil
}
