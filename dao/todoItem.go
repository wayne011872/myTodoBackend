package dao

import(
	"time"
)
type TodoItem struct {
	ID          int64  
	Title       string 
	Detail      string 
	Completed   bool 
	StartTime   time.Time
	EndTime     time.Time
	CreateTime  time.Time `json:"createtime,omitempty"`
	UpdateTime  time.Time `json:"updatetime,omitempty"`
}