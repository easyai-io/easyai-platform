package schema

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

// EventReasonSummary 聚合job‘s tasks关联的同一类型k8s event
type EventReasonSummary struct {
	Reason         string    `json:"reason" db:"reason"`        // event reason
	Type           string    `json:"type" db:"event_type"`      // event type
	FirstTimestamp time.Time `json:"timestamp" db:"event_time"` // 最早的出现时间
	Total          int       `json:"total" db:"total"`
	Count          int       `json:"count" db:"count"`
}

// EasyaiEvent 自定义事件
type EasyaiEvent struct {
	K8SEvent       corev1.Event `json:"-"`
	Name           string       `json:"name"`
	Source         string       `json:"source"`
	Message        string       `json:"message"`
	Reason         string       `json:"reason"`
	Type           string       `json:"type"`
	ObjectName     string       `json:"object_name"`
	FirstTimestamp time.Time    `json:"timestamp"`
	LastTimestamp  time.Time    `json:"last_timestamp"`

	JobID uint32 `json:"job_id"`
	Count int    `json:"count"`
	Owner string `json:"owner"`
}
