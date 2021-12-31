package client

import (
	"bytes"
	"fmt"
	"github.com/ZYallers/rpcx-framework/helper"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/smallnest/rpcx/codec"
	"math/rand"
	"strconv"
	"time"
)

func HttpInvoke(serviceName, serviceAddr, serviceMethod string, args map[string]interface{}, other ...interface{}) (interface{}, error) {
	req := tool.NewRequest(fmt.Sprintf("http://%s", serviceAddr))
	req.SetHeaders(map[string]string{
		"X-RPCX-Version":       "1.6.11",
		"X-RPCX-MessageID":     strconv.Itoa(rand.Int()),
		"X-RPCX-MesssageType":  "0",
		"X-RPCX-SerializeType": "3",
		"X-RPCX-ServicePath":   serviceName,
		"X-RPCX-ServiceMethod": serviceMethod,
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
	if err := cc.Decode(helper.String2Bytes(resp.Body), &reply); err != nil {
		return nil, err
	}
	return reply, nil
}
