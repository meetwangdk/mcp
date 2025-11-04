package model

type EavesdroppingMeta struct {
	ClusterName    string      `json:"clusterName"`
	LastReadTime   interface{} `json:"LastReadTime"`
	LivePodLatency interface{} `json:"Latency"`
}

func (*EavesdroppingMeta) TableName() string {
	return "eavesdripping_meta"
}
func (s *EavesdroppingMeta) EsTableName() string {
	return "eavesdripping_meta"
}

func (s *EavesdroppingMeta) TypeName() string {
	return "meta"
}
