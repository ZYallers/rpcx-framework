package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/ZYallers/golib/utils/curl"
	"github.com/ZYallers/rpcx-framework/consts"
	"github.com/smallnest/rpcx/log"
)

const (
	timeout   = 3 * time.Second
	uriPrefix = "https://oapi.dingtalk.com/robot/send?access_token="
)

var (
	header  = map[string]string{"Content-Type": "application/json;charset=utf-8"}
	message *Message
)

type Message struct {
	ErrorToken    string
	GracefulToken string
	SqlToken      string
	Mode          string
	Name          string
	Addr          string
	Hostname      string
	SystemIP      string
	PublicIP      string
}

func InitMessage(m *Message)    { message = m }
func GetMessage() *Message      { return message }
func (s *Message) Open() bool   { return s != nil }
func (s *Message) Always() bool { return s != nil && s.Mode == consts.DevelopMode }
func (s *Message) Push(msg string) {
	if s != nil {
		s.Send(s.SqlToken, msg, true)
	}
}

func (s *Message) Graceful(msg interface{}, isAtAll bool, logType ...interface{}) {
	if s != nil {
		s.Send(s.GracefulToken, msg, append([]interface{}{"", isAtAll}, logType...)...)
	}
}

func (s *Message) Error(msg interface{}, stack string, isAtAll bool, logType ...interface{}) {
	if s != nil {
		s.Send(s.ErrorToken, msg, append([]interface{}{stack, isAtAll}, logType...)...)
	}
}

func (s *Message) Send(token string, msg interface{}, options ...interface{}) {
	defer func() { recover() }()
	if s == nil {
		return
	}
	title := fmt.Sprintf("%v", msg)
	if token == "" || title == "" {
		return
	}
	text := []string{
		title + "\n---------------------------",
		"Env: " + s.Mode,
		"Name: " + s.Name,
		"Addr: " + s.Addr,
		"HostName: " + s.Hostname,
		"Time: " + time.Now().Format("2006/01/02 15:04:05.000"),
		"SystemIP: " + s.SystemIP,
		"PublicIP: " + s.PublicIP,
	}
	if len(options) > 0 {
		if stack, ok := options[0].(string); ok && stack != "" {
			text = append(text, "\nStack:\n"+stack)
		}
	}
	var isAtAll bool
	if s.Mode == consts.ProduceMode && len(options) > 1 {
		if val, ok := options[1].(bool); ok {
			isAtAll = val
		}
	}
	data := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": strings.Join(text, "\n") + "\n"},
		"at":      map[string]interface{}{"isAtAll": isAtAll},
	}
	_, _ = curl.NewRequest(uriPrefix + token).SetHeaders(header).SetPostData(data).SetTimeOut(timeout).Post()
	if len(options) > 2 {
		if logType, ok := options[2].(string); ok && logType != "" {
			switch strings.ToLower(logType) {
			case "debug":
				log.Debug(title)
			case "info":
				log.Info(title)
			case "warn":
				log.Warn(title)
			case "error":
				log.Error(title)
			case "fatal":
				log.Fatal(title)
			case "panic":
				log.Panic(title)
			}
		}
	}
}
