package model

import "time"

type SloOptions struct {
	From               time.Time // range query
	To                 time.Time // range query
	Result             string
	Cluster            string
	BizName            string
	NodeName           string
	Count              string
	Filter             string
	DeliveryDuration   string
	Type               string // create 或者 delete, 默认是 create
	DeliveryStatus     string // FAIL/KILL/ALL/SUCCESS
	SloTime            string // 20s/30m0s/10m0s
	Env                string // prod, test
	NativeQuery        string // native query
	UserError          string // 用户错误
	Namespace          string // 增加namespace查询
	TimeField          string // 表示根据哪个时间字段进行查询
	RemovingUserErrors bool   // 移除用户错误标识
}
type NodeParams struct {
	NodeUid  string
	NodeIp   string
	NodeName string
}
type PodParams struct {
	Name     string
	Uid      string
	Hostname string
	Podip    string
}
