package types

import (
	"time"
)

type ReplyService struct {
	Name     string     `json:"name,omitempty"`
	Hostname string     `json:"hostname,omitempty"`
	Ip       string     `json:"ip,omitempty"`
	Addr     string     `json:"addr,omitempty"`
	Runtime  string     `json:"runtime,omitempty"`
	Start    *time.Time `json:"start,omitempty"`
	End      *time.Time `json:"end,omitempty"`
}

type Record struct {
	Type      string      `json:"type,omitempty"`
	TableName string      `json:"table_name,omitempty"`
	DataId    interface{} `json:"data_id,omitempty"`
	Intro     string      `json:"intro,omitempty"`
}

type Reply struct {
	Code    int           `json:"code"`
	Msg     string        `json:"msg"`
	Data    interface{}   `json:"data"`
	Record  *Record       `json:"record,omitempty"`
	Service *ReplyService `json:"service,omitempty"`
}
