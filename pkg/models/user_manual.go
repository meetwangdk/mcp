package models

import (
	"time"
)

type UserManualPushRecord struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     string    `gorm:"index;not null" json:"user_id"` // 用户唯一标识
	LastPushAt time.Time `gorm:"not null" json:"last_push_at"`  // 最后推送时间
}

func (s *UserManualPushRecord) TableName() string {
	return "user_manual_push_record"
}

func (s *UserManualPushRecord) EsTableName() string {
	return "user_manual_push_record"
}

func (s *UserManualPushRecord) TypeName() string {
	return "_doc"
}
