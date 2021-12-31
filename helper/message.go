package helper

import (
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/smallnest/rpcx/log"
	"github.com/spf13/viper"
	"github.com/ZYallers/rpcx-framework/env"
	"os"
	"strings"
	"time"
)

const (
	robotLinkPrefix           = "https://oapi.dingtalk.com/robot/send?access_token="
	defaultErrorRobotToken    = "xxxxxx"
	defaultGracefulRobotToken = "xxxxxx"
)

var (
	header = map[string]string{"Content-Type": "application/json;charset=utf-8"}
)

func readConfig(key string) string {
	return viper.GetString(key)
}

func serviceName() string {
	name := readConfig("service.name")
	if name == "" {
		return "unknown"
	}
	return name
}

func serviceAddr() string {
	addr := readConfig("service.addr")
	if addr == "" {
		return "unknown"
	}
	return addr
}

func serviceMode() string {
	mode := "unknown"
	modeKey := readConfig("global.modeKey")
	if modeKey == "" {
		return mode
	}
	if val := os.Getenv(modeKey); val != "" {
		mode = val
	}
	return mode
}

func hostname() string {
	hostname, _ := os.Hostname()
	if hostname == "" {
		return "unknown"
	}
	return hostname
}

func errorRobotUrl() string {
	token := readConfig("service.errorRobotToken")
	if token != "" {
		return robotLinkPrefix + token
	}
	return robotLinkPrefix + defaultErrorRobotToken
}

func gracefulRobotUrl() string {
	token := readConfig("service.gracefulRobotToken")
	if token != "" {
		return robotLinkPrefix + token
	}
	return robotLinkPrefix + defaultGracefulRobotToken
}

func ContextMessage(msg string, stack string, isAtAll bool, logType ...interface{}) {
	mode := serviceMode()
	text := []string{
		msg + "\n---------------------------",
		"Env: " + mode,
		"Name: " + serviceName(),
		"Addr: " + serviceAddr(),
		"HostName: " + hostname(),
		"Time: " + time.Now().Format("2006/01/02 15:04:05.000"),
		"SystemIP: " + tool.SystemIP(),
		"PublicIP: " + tool.PublicIP(),
	}
	if stack != "" {
		text = append(text, "\nStack:\n"+stack)
	}
	if mode != env.ProduceMode {
		isAtAll = false
	}
	postData := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": strings.Join(text, "\n") + "\n"},
		"at":      map[string]interface{}{"isAtAll": isAtAll},
	}
	if resp, err := tool.NewRequest(errorRobotUrl()).SetHeaders(header).SetPostData(postData).
		SetTimeOut(3 * time.Second).Post(); err != nil {
		log.Errorf("context message error: %v, resp: %v, msg: %s", err, resp, msg)
	}
	if len(logType) == 1 {
		writeLog(msg, logType[0].(string))
	}
}

func SimpleMessage(msg string, isAtAll bool, logType ...interface{}) {
	mode := serviceMode()
	text := []string{
		msg + "\n---------------------------",
		"Env: " + mode,
		"Name: " + serviceName(),
		"Addr: " + serviceAddr(),
		"HostName: " + hostname(),
		"Time: " + time.Now().Format("2006/01/02 15:04:05.000"),
		"SystemIP: " + tool.SystemIP(),
		"PublicIP: " + tool.PublicIP(),
	}
	if mode != env.ProduceMode {
		isAtAll = false
	}
	postData := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": strings.Join(text, "\n") + "\n"},
		"at":      map[string]interface{}{"isAtAll": isAtAll},
	}
	if resp, err := tool.NewRequest(gracefulRobotUrl()).SetHeaders(header).SetPostData(postData).
		SetTimeOut(3 * time.Second).Post(); err != nil {
		log.Errorf("simple message error: %v, resp: %v, msg: %s", err, resp, msg)
	}
	if len(logType) == 1 {
		writeLog(msg, logType[0].(string))
	}
}

func writeLog(msg, logType string) {
	switch strings.ToLower(logType) {
	case "debug":
		log.Debug(msg)
	case "info":
		log.Info(msg)
	case "warn":
		log.Warn(msg)
	case "error":
		log.Error(msg)
	case "fatal":
		log.Fatal(msg)
	case "panic":
		log.Panic(msg)
	}
}
