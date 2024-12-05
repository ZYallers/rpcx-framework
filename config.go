package framework

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/ZYallers/golib/funcs/nets"
	"github.com/ZYallers/rpcx-framework/consts"
	"github.com/ZYallers/rpcx-framework/types"
	"github.com/spf13/viper"
)

var (
	serviceMode               string
	serviceName               string
	serviceAddr               string
	serviceHostname           string
	serviceLogDir             string
	serviceVersion            string
	serviceVersionKey         string
	serviceErrorRobotToken    string
	serviceGracefulRobotToken string
	serviceSqlRobotToken      string
	systemIP                  string
	publicIP                  string
	serviceDiscovery          *types.Discovery
)

func ReadInConfig(args ...string) {
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
		panic(fmt.Errorf("read config file '%s/%s.%s' error: %s", configPath, configName, configType, err))
	}
}

func ServiceMode() string {
	if serviceMode == "" {
		if modeKey := viper.GetString("global.modeKey"); modeKey == "" {
			return consts.DevelopMode
		} else {
			if s := os.Getenv(modeKey); s == "" {
				return consts.DevelopMode
			} else {
				serviceMode = s
			}
		}
	}
	return serviceMode
}

func ServiceName() string {
	if serviceName == "" {
		if s := viper.GetString("service.name"); s != "" {
			serviceName = s
		}
	}
	return serviceName
}

func ServiceAddr() string {
	if serviceAddr == "" {
		if s := viper.GetString("service.addr"); s != "" {
			if ip := SystemIP(); ip != "" {
				serviceAddr = strings.Replace(s, "0.0.0.0", ip, 1)
			}
		}
	}
	return serviceAddr
}

func ServiceHostname() string {
	if serviceHostname == "" {
		if s, _ := os.Hostname(); s != "" {
			serviceHostname = strings.ToLower(s)
		}
	}
	return serviceHostname
}

func ServiceLogDir() string {
	if serviceLogDir == "" {
		if s := viper.GetString("service.logDir"); s != "" {
			serviceLogDir = s
		}
	}
	return serviceLogDir
}

func ServiceVersion() string {
	if serviceVersion == "" {
		if s := viper.GetString("service.version"); s != "" {
			serviceVersion = s
		}
	}
	return serviceVersion
}

func ServiceVersionKey() string {
	if serviceVersionKey == "" {
		if s := viper.GetString("service.versionKey"); s != "" {
			serviceVersionKey = s
		}
	}
	return serviceVersionKey
}

func ServiceErrorRobotToken() string {
	if serviceErrorRobotToken == "" {
		if s := viper.GetString("service.errorRobotToken"); s != "" {
			serviceErrorRobotToken = s
		}
	}
	return serviceErrorRobotToken
}

func ServiceGracefulRobotToken() string {
	if serviceGracefulRobotToken == "" {
		if s := viper.GetString("service.gracefulRobotToken"); s != "" {
			serviceGracefulRobotToken = s
		}
	}
	return serviceGracefulRobotToken
}

func ServiceSqlRobotToken() string {
	if serviceSqlRobotToken == "" {
		if s := viper.GetString("service.sqlRobotToken"); s != "" {
			serviceSqlRobotToken = s
		}
	}
	return serviceSqlRobotToken
}

func ServiceDiscovery() *types.Discovery {
	if serviceDiscovery != nil {
		return serviceDiscovery
	}

	mode := ServiceMode()
	basePath := viper.GetString(fmt.Sprintf("service.etcd.%s.basePath", mode))
	if basePath == "" {
		return nil
	}

	addr := viper.GetString(fmt.Sprintf("service.etcd.%s.addr", mode))
	if len(addr) > 0 && addr[0:1] == "$" {
		addr = os.Getenv(addr[1:])
	}
	if addr == "" {
		return nil
	}

	if mode == consts.DevelopMode {
		if hostname := ServiceHostname(); hostname != viper.GetString("global.server.development.hostname") {
			if developerDockerHostname := os.Getenv("developer_docker_hostname"); developerDockerHostname != "" {
				hostname = developerDockerHostname
			}
			basePath = strings.Replace(basePath, consts.DevelopMode, "developer@"+hostname, 1)
			addr = strings.Replace(addr, "127.0.0.1", viper.GetString("global.server.development.ip"), 1)
		}
	}

	interval := viper.GetInt64(fmt.Sprintf("service.etcd.%s.updateInterval", mode))
	if interval <= 0 {
		interval = 30
	}

	serviceDiscovery = &types.Discovery{
		BasePath:       basePath,
		Addr:           strings.Split(addr, ","),
		UpdateInterval: time.Duration(interval) * time.Second,
	}

	return serviceDiscovery
}

func SystemIP() string {
	if systemIP == "" {
		if s := nets.SystemIP(); s != "" && s != "unknown" {
			systemIP = s
		}
	}
	return systemIP
}

func PublicIP() string {
	if publicIP == "" {
		if s := nets.PublicIP(); s != "" && s != "unknown" {
			publicIP = s
		}
	}
	return publicIP
}
