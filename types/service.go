package types

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ZYallers/golib/funcs/conv"
	"github.com/ZYallers/golib/funcs/php"
	"github.com/ZYallers/golib/utils/json"
	"github.com/smallnest/rpcx/server"
	"github.com/spf13/viper"
)

type M map[string]interface{}

type IService interface {
	Construct(service interface{}, ctx context.Context, args map[string]interface{}, reply *interface{})
	SignCheck() bool
	LoginCheck(values ...string) bool
}

type Service struct {
	debug   bool
	ctx     context.Context
	service *Rpc
	args    map[string]interface{}
	reply   *interface{}
}

func (s *Service) getServiceConfig(key string) interface{} {
	return viper.Get("service." + key)
}

func (s *Service) Construct(sr interface{}, ctx context.Context, args map[string]interface{}, reply *interface{}) {
	s.service = sr.(*Rpc)
	s.ctx = ctx
	s.args = args
	s.reply = reply
	rep := &Reply{}
	if debug := s.GetString("debug"); debug == conv.ToString(s.getServiceConfig("debugValue")) {
		s.debug = true
		now := time.Now()
		rep.Service = &ReplyService{
			Name:     s.service.Name,
			Hostname: s.service.HostName,
			Ip:       s.service.SystemIP,
			Addr:     s.service.Addr,
			Start:    &now,
		}
	}
	*s.reply = rep
}

func (s *Service) GetArgs(key string, defaultValue ...interface{}) interface{} {
	if v, ok := s.args[key]; ok {
		return v
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (s *Service) GetString(key string, defaultValue ...string) string {
	if v, ok := s.args[key]; ok {
		if s, err := conv.ToStringE(v); err != nil {
			return fmt.Sprint(v)
		} else {
			return s
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (s *Service) GetInt(key string, defaultValue ...int) int {
	if v, ok := s.args[key]; ok {
		if i, err := conv.ToIntE(v); err != nil {
			i, _ = strconv.Atoi(fmt.Sprint(v))
			return i
		} else {
			return i
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (s *Service) GetInt64(key string, defaultValue ...int64) int64 {
	if v, ok := s.args[key]; ok {
		if i, err := conv.ToInt64E(v); err != nil {
			i, _ = strconv.ParseInt(fmt.Sprint(v), 10, 64)
			return i
		} else {
			return i
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (s *Service) Finish(data interface{}) error {
	resp := (*s.reply).(*Reply)
	resp.Data = data
	if s.debug {
		now := time.Now()
		resp.Service.End = &now
		resp.Service.Runtime = time.Since(*resp.Service.Start).String()
	}
	*s.reply = resp
	return nil
}

func (s *Service) Record(r Record) {
	resp := (*s.reply).(*Reply)
	resp.Record = &r
	*s.reply = resp
}

func (s *Service) Json(a ...interface{}) error {
	rep := (*s.reply).(*Reply)
	rep.Code = http.StatusOK
	rep.Msg = ""
	rep.Data = struct{}{}
	al := len(a)
	if al > 0 {
		rep.Code = conv.ToInt(a[0])
	}
	if al > 1 {
		rep.Msg = conv.ToString(a[1])
	}
	if al > 2 {
		rep.Data = a[2]
	}
	if s.debug {
		now := time.Now()
		rep.Service.End = &now
		rep.Service.Runtime = time.Since(*rep.Service.Start).String()
	}
	if s.ctx.Value(server.HttpConnContextKey) == nil {
		bte, err := json.Marshal(rep)
		if err != nil {
			return s.Json(http.StatusInternalServerError, err)
		}
		*s.reply = bte
	} else {
		*s.reply = rep
	}
	return nil
}

func (s *Service) WriteString(str string) error {
	*s.reply = str
	return nil
}

func (s *Service) SignCheck() bool {
	sign := s.GetString(conv.ToString(s.getServiceConfig("signKey")))
	if sign == "" {
		return false
	}
	utime := s.GetString(conv.ToString(s.getServiceConfig("utimeKey")))
	if utime == "" {
		return false
	}
	timestamp, err := strconv.ParseInt(utime, 10, 0)
	if err != nil || timestamp <= 0 {
		return false
	}
	if time.Now().Unix()-timestamp > conv.ToInt64(s.getServiceConfig("signExpire")) {
		return false
	}
	hash := md5.New()
	hash.Write([]byte(utime + conv.ToString(s.getServiceConfig("signSecret"))))
	realSign := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(hash.Sum(nil))))
	if sign == realSign {
		return true
	}
	return false
}

func (s *Service) LoggedUserData(key ...string) map[string]interface{} {
	var token string
	switch len(key) {
	case 1:
		token = key[0]
	default:
		token = s.GetString(conv.ToString(s.getServiceConfig("tokenKey")))
	}
	if s.service.SessionFunc == nil {
		return nil
	}
	session := s.service.SessionFunc()
	if session == nil {
		return nil
	}
	if str, _ := session.Get(conv.ToString(s.getServiceConfig("sessionKeyPrefix")) + token).Result(); str != "" {
		return php.Unserialize(str)
	}
	return nil
}

func (s *Service) LoginCheck(key ...string) bool {
	if vars := s.LoggedUserData(key...); vars != nil {
		return true
	}
	return false
}

func (s *Service) LoggedUserId(key ...string) int {
	vars := s.LoggedUserData(key...)
	if vars == nil {
		return 0
	}
	if data, ok := vars["userinfo"].(map[string]interface{}); ok {
		if str, ok := data["userid"].(string); ok && str != "" {
			return conv.ToInt(str)
		}
	}
	return 0
}
