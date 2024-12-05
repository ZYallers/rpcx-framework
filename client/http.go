package client

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/ZYallers/golib/utils/curl"
	"github.com/smallnest/rpcx/codec"
)

func HttpInvoke(serviceName, serviceAddr, serviceMethod string, args map[string]interface{}, other ...interface{}) (interface{}, error) {
	req := curl.NewRequest("http://" + serviceAddr)
	req.SetHeaders(map[string]string{
		"X-RPCX-Version":       "1.6.11",
		"X-RPCX-MesssageType":  "0",
		"X-RPCX-SerializeType": "3",
		"X-RPCX-ServicePath":   serviceName,
		"X-RPCX-ServiceMethod": serviceMethod,
		"X-RPCX-MessageID":     strconv.Itoa(rand.Int()),
	})
	cc := &codec.MsgpackCodec{}
	data, _ := cc.Encode(args)
	req.SetBody(bytes.NewReader(data))
	if len(other) > 0 {
		req.SetTimeOut(other[0].(time.Duration))
	}
	resp, err := req.Post()
	if err != nil {
		return nil, err
	}
	status := resp.Raw.Status
	statusCode := resp.Raw.StatusCode
	errMsg := resp.Raw.Header.Get("X-Rpcx-Errormessage")
	if statusCode != 200 {
		return nil, fmt.Errorf("response error: code:%d, status:%s, message:%s", statusCode, status, errMsg)
	}
	if resp.Body == "" {
		return nil, nil
	}
	var reply interface{}
	if err := cc.Decode([]byte(resp.Body), &reply); err != nil {
		return nil, err
	}
	return reply, nil
}
