package schema

import "time"

// JobDemo is a demo for job
type JobDemo struct {
	ID          uint32            `json:"id" yaml:"id"`
	Name        string            `json:"name" yaml:"name"`
	Framework   JobFramework      `json:"framework" yaml:"framework"`
	Description string            `json:"description" yaml:"description"`
	Owner       string            `json:"owner" yaml:"owner"`
	DocLinks    map[string]string `json:"doc_links" yaml:"doc_links"`
	JobLinks    map[uint32]string `json:"job_links" yaml:"job_links"`
	CreatedAt   time.Time         `json:"created_at" yaml:"created_at"`
	ModifiedAt  time.Time         `json:"modified_at" yaml:"modified_at"`
}
