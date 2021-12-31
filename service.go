package framework

import (
	"fmt"
	"github.com/ZYallers/rpcx-framework/env"
	"github.com/ZYallers/rpcx-framework/helper"
	"github.com/ZYallers/rpcx-framework/service"
	"github.com/ZYallers/rpcx-framework/util/zap"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/smallnest/rpcx/log"
	"github.com/spf13/viper"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

var rpcxService *service.RPCXService

func LoadConfig(args ...string) {
	relativePath, configName, configType := ".", "service", "json"
	argsLen := len(args)
	if argsLen > 0 {
		relativePath = args[0]
	}
	if argsLen > 1 {
		configName = args[1]
	}
	if argsLen > 2 {
		configType = args[2]
	}
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	_, filePath, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filePath), relativePath)
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	if err := viper.ReadInConfig(); err != nil {
		helper.SimpleMessage(fmt.Sprintf("read config file error: %s", err), true, "panic")
	}
}

func ReadConfig(key string) string {
	return viper.GetString(key)
}

func ServiceName() string {
	if rpcxService != nil {
		return rpcxService.Name
	}
	return ReadConfig("service.name")
}

func ServiceMode() string {
	if rpcxService != nil {
		return rpcxService.Env
	}
	mode := env.DevelopMode
	modeKey := ReadConfig("global.modeKey")
	if modeKey == "" {
		return mode
	}
	if val := os.Getenv(modeKey); val != "" {
		mode = val
	}
	return mode
}

func ServiceHostName() string {
	if rpcxService != nil {
		return rpcxService.HostName
	}
	hostname, _ := os.Hostname()
	if hostname != "" {
		hostname = strings.ToLower(hostname)
	}
	return hostname
}

func DiscoveryBasePath() string {
	if rpcxService != nil {
		return rpcxService.Etcd.BasePath
	}
	runMode := ServiceMode()
	basePath := ReadConfig(fmt.Sprintf("service.etcd.%s.basePath", runMode))
	if basePath == "" {
		return ""
	}
	hostname := ServiceHostName()
	if runMode == env.DevelopMode && hostname != ReadConfig("global.server.development.hostname") {
		basePath = strings.Replace(basePath, env.DevelopMode, "developer@"+hostname, 1)
	}
	return basePath
}

func DiscoveryAddress() []string {
	if rpcxService != nil {
		return rpcxService.Etcd.Addr
	}
	runMode := ServiceMode()
	addr := ReadConfig(fmt.Sprintf("service.etcd.%s.addr", runMode))
	if addr == "" {
		return nil
	}
	hostname := ServiceHostName()
	if runMode == env.DevelopMode && hostname != ReadConfig("global.server.development.hostname") {
		developmentServerIP := ReadConfig("global.server.development.ip")
		addr = strings.Replace(addr, "127.0.0.1", developmentServerIP, 1)
	}
	return strings.Split(addr, ",")
}

func DiscoveryUpdateInterval() time.Duration {
	if rpcxService != nil {
		return rpcxService.Etcd.UpdateInterval
	}
	runMode := ServiceMode()
	interval := viper.GetInt64(fmt.Sprintf("service.etcd.%s.updateInterval", runMode))
	if interval <= 0 {
		return time.Duration(30) * time.Second
	}
	return time.Duration(interval) * time.Second
}

func NewService() *service.RPCXService {
	defer func() {
		if err := recover(); err != nil {
			helper.SimpleMessage(fmt.Sprintf("%v", err), true, "panic")
		}
	}()

	systemIP := tool.SystemIP()
	if systemIP == "unknown" || systemIP == "" {
		panic("system ip is unknown or empty")
	}

	hostname := ServiceHostName()
	if hostname == "" {
		panic("system hostname is empty")
	}

	serviceName := ServiceName()
	if serviceName == "" {
		panic("service name is empty")
	}

	serviceLogDir := ReadConfig("service.logDir")
	if serviceLogDir == "" {
		panic("service log dir is empty")
	}

	zap.SetLoggerDir(serviceLogDir)
	log.SetLogger(service.NewLogger(serviceName))

	discoveryBasePath := DiscoveryBasePath()
	if discoveryBasePath == "" {
		panic("discovery base path is empty")
	}

	discoveryAddress := DiscoveryAddress()
	if len(discoveryAddress) == 0 {
		panic("discovery address is empty")
	}

	rpcxService = &service.RPCXService{
		Env:                ServiceMode(),
		Name:               serviceName,
		HostName:           hostname,
		SystemIP:           systemIP,
		LogDir:             serviceLogDir,
		Addr:               strings.Replace(ReadConfig("service.addr"), "0.0.0.0", systemIP, 1),
		Version:            ReadConfig("service.version"),
		VersionKey:         ReadConfig("service.versionKey"),
		ErrorRobotToken:    ReadConfig("service.errorRobotToken"),
		GracefulRobotToken: ReadConfig("service.gracefulRobotToken"),
		Etcd: &service.Discovery{
			BasePath:       discoveryBasePath,
			Addr:           discoveryAddress,
			UpdateInterval: DiscoveryUpdateInterval(),
		},
	}

	return rpcxService
}
