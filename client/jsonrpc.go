package client

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/ZYallers/golib/utils/curl"
	"github.com/ZYallers/golib/utils/json"
)

type jsonRpc struct {
	Id      int                    `json:"id"`
	Jsonrpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	Result  interface{}            `json:"result"`
	Error   struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
}

func JsonRpc2(serviceName, serviceAddr, serviceMethod string, args map[string]interface{}, options ...interface{}) (interface{}, error) {
	req := curl.NewRequest(fmt.Sprintf("http://%s", serviceAddr))
	req.SetHeaders(map[string]string{"X-JSONRPC-2.0": "true"})
	data := jsonRpc{
		Jsonrpc: "2.0",
		Id:      rand.Int(),
		Method:  fmt.Sprintf("%s.%s", serviceName, serviceMethod),
		Params:  args,
	}
	b, _ := json.Marshal(data)
	req.SetBody(bytes.NewReader(b))
	if len(options) > 0 {
		req.SetTimeOut(options[0].(time.Duration))
	}
	resp, err := req.Post()
	if err != nil {
		return nil, err
	}
	if resp.Body == "" {
		return nil, nil
	}
	var res jsonRpc
	if err := json.Unmarshal([]byte(resp.Body), &res); err != nil {
		return nil, err
	}
	if res.Error.Message != "" {
		return nil, fmt.Errorf("jsonRpc2 error: code:%d, message:%s, data:%v", res.Error.Code, res.Error.Message, res.Error.Data)
	}
	return res.Result, nil
}
