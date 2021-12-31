package mtsc

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ZYallers/rpcx-framework/define"
	"github.com/ZYallers/rpcx-framework/env"
	"github.com/ZYallers/rpcx-framework/service"
	"github.com/ZYallers/zgin/libraries/json"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/smallnest/rpcx/server"
	"net/http"
	"strconv"
	"time"
)

const (
	debugValue = "HXSAPP2021"
	tokenKey   = "sess_token"
	signKey    = "sign"
	devSignKey = "hxs-rpcx-dev"
	utimeKey   = "utime"
)

type Service struct {
	debug   bool
	ctx     context.Context
	service *service.RPCXService
	args    map[string]interface{}
	reply   *interface{}
}

func (s *Service) Construct(sr interface{}, ctx context.Context, args map[string]interface{}, reply *interface{}) {
	s.service = sr.(*service.RPCXService)
	s.ctx = ctx
	s.args = args
	s.reply = reply
	rep := &define.Reply{}
	if debug := s.GetArgs("debug"); debug == debugValue {
		s.debug = true
		now := time.Now()
		rep.Service = &define.ReplyService{
			Name:     s.service.Name,
			Hostname: s.service.HostName,
			Ip:       s.service.SystemIP,
			Addr:     s.service.Addr,
			Start:    &now,
		}
	}
	*s.reply = rep
}

//  GetArgs 获取客户端传参值
//  @receiver s *Service
//  @author Cloud|2021-12-02 16:18:05
//  @param key string ...
//  @param defaultValue ...interface{} ...
//  @return interface{} ...
func (s *Service) GetArgs(key string, defaultValue ...interface{}) interface{} {
	if val, exist := s.args[key]; exist {
		return val
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

//  GetString 获取客户端字符串传值
//  @receiver s *Service
//  @author Cloud|2021-12-23 11:41:56
//  @param key string ...
//  @param defaultValue ...string ...
//  @return string ...
func (s *Service) GetString(key string, defaultValue ...string) string {
	if val, exist := s.args[key]; exist {
		var res string
		switch v := val.(type) {
		case int:
			res = strconv.Itoa(v)
		case int8:
			res = strconv.Itoa(int(v))
		default:
			res = fmt.Sprintf("%v", v)
		}
		return res
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

//  GetInt .获取客户端int传值
//  @receiver s *Service
//  @author Cloud|2021-12-23 11:42:24
//  @param key string ...
//  @param defaultValue ...int ...
//  @return int ...
func (s *Service) GetInt(key string, defaultValue ...int) int {
	if val, exist := s.args[key]; exist {
		var res int
		switch v := val.(type) {
		case string:
			res, _ = strconv.Atoi(v)
		case int8:
			res = int(v)
		case int16:
			res = int(v)
		case int32:
			res = int(v)
		default:
			res, _ = strconv.Atoi(fmt.Sprintf("%d", v))
		}
		return res
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

//  GetInt64 ...
//  @receiver s *Service
//  @author Cloud|2021-12-23 11:55:53
//  @param key string ...
//  @param defaultValue ...int64 ...
//  @return int64 ...
func (s *Service) GetInt64(key string, defaultValue ...int64) int64 {
	if val, exist := s.args[key]; exist {
		var str string
		switch v := val.(type) {
		case string:
			str = v
		default:
			str = fmt.Sprintf("%s", v)
		}
		if str == "" {
			return 0
		}
		i, err := strconv.ParseInt(str, 10, 0)
		if err != nil {
			return 0
		}
		return i
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

//  Finish ...
//  @receiver s *Service
//  @author Cloud|2021-12-23 11:42:27
//  @param data interface{} ...
//  @return error ...
func (s *Service) Finish(data interface{}) error {
	resp := (*s.reply).(*define.Reply)
	resp.Data = data
	if s.debug {
		now := time.Now()
		resp.Service.End = &now
		resp.Service.Runtime = time.Since(*resp.Service.Start).String()
	}
	*s.reply = resp
	return nil
}

//  Record ...
//  @receiver s *Service
//  @author Cloud|2021-12-07 15:31:39
//  @param r define.Record ...
func (s *Service) Record(r define.Record) {
	resp := (*s.reply).(*define.Reply)
	resp.Record = &r
	*s.reply = resp
}

//  Json ...
//  @receiver s *Service
//  @author Cloud|2021-12-07 14:12:03
//  @param args ...interface{} ...
//  @return error ...
func (s *Service) Json(args ...interface{}) error {
	rep := (*s.reply).(*define.Reply)
	if len(args) == 0 {
		rep.Code = http.StatusOK
	}
	if len(args) > 0 {
		rep.Code = args[0].(int)
	}
	if len(args) > 1 {
		switch value := args[1].(type) {
		case error:
			rep.Msg = value.Error()
		case string:
			rep.Msg = value
		default:
			rep.Msg = fmt.Sprintf("%v", value)
		}
	}
	if len(args) > 2 {
		rep.Data = args[2]
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

//  SignCheck APP签名验证
//  @receiver s *Service
//  @author Cloud|2021-12-02 16:57:44
//  @return bool ...
func (s *Service) SignCheck() bool {
	sign := s.GetString(signKey)
	if sign == "" {
		return false
	}
	if s.service.Env == env.DevelopMode && sign == devSignKey {
		return true
	}
	utime := s.GetString(utimeKey)
	if utime == "" {
		return false
	}
	timestamp, err := strconv.ParseInt(utime, 10, 0)
	if err != nil || timestamp <= 0 {
		return false
	}
	if time.Now().Unix()-timestamp > env.SignTimeExpiration {
		return false
	}
	hash := md5.New()
	hash.Write([]byte(utime + env.TokenKey))
	realSign := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(hash.Sum(nil))))
	if sign == realSign {
		return true
	}
	return false
}

//  LoggedUserData 获取APP登录用户数据
//  @receiver s *Service
//  @author Cloud|2021-12-02 17:08:20
//  @param key ...string ...
//  @return map[string]interface{} ...
func (s *Service) LoggedUserData(key ...string) map[string]interface{} {
	var token string
	switch len(key) {
	case 1:
		token = key[0]
	default:
		token = s.GetString(tokenKey)
	}
	client := s.service.SessionClient()
	if client == nil {
		return nil
	}
	if str, _ := client.Get("ci_session:" + token).Result(); str == "" {
		return nil
	} else {
		return tool.PhpUnserialize(str)
	}
}

//  LoginCheck APP登录检查
//  @receiver s *Service
//  @author Cloud|2021-12-02 17:07:53
//  @param key ...string ...
//  @return bool ...
func (s *Service) LoginCheck(key ...string) bool {
	if vars := s.LoggedUserData(key...); vars != nil {
		return true
	}
	return false
}

//  LoggedUserId 获取APP登陆用户的user_id
//  @receiver s *Service
//  @author Cloud|2021-12-22 18:38:51
//  @param key ...string ...
//  @return int ...
func (s *Service) LoggedUserId(key ...string) int {
	vars := s.LoggedUserData(key...)
	if vars == nil {
		return 0
	}
	if data, ok := vars["userinfo"].(map[string]interface{}); ok {
		if str, ok := data["userid"].(string); ok && str != "" {
			userId, _ := strconv.Atoi(str)
			return userId
		}
	}
	return 0
}
