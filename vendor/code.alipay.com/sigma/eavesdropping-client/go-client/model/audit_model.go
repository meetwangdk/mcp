package model

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sync"
	"time"

	k8saudit "k8s.io/apiserver/pkg/apis/audit"
)

type Audit struct {
	AuditId           string         `json:"auditId,omitempty" gorm:"audit_id"`                  // 主键
	SqlStageTimestamp time.Time      `json:"sqlStageTimestamp,omitempty" gorm:"stage_timestamp"` // stage_timestamp
	Cluster           string         `json:"cluster,omitempty" gorm:"cluster"`                   // cluster
	Namespace         string         `json:"namespace,omitempty" gorm:"namespace"`               // namespace
	Resource          string         `json:"resource,omitempty" gorm:"resource"`                 // resource
	Content           string         `json:"content,omitempty" gorm:"content"`                   // content
	AuditLog          k8saudit.Event `json:"auditLog,omitempty" gorm:"-"`
}

func (*Audit) TableName() string {
	return "audit"
}
func (s *Audit) EsTableName() string {
	return "audit_*"
}

func (s *Audit) TypeName() string {
	return "doc"
}

type AuditEvent struct {
	sync.WaitGroup
	*k8saudit.Event
	ResponseRuntimeObj runtime.Object `json:"responseRuntimeObj,omitempty"`
	RequestRuntimeObj  runtime.Object `json:"requestRuntimeObj,omitempty"`
}

func (*AuditEvent) TableName() string {
	return "audit"
}
func (s *AuditEvent) EsTableName() string {
	return "audit_*"
}

func (s *AuditEvent) TypeName() string {
	return "doc"
}
