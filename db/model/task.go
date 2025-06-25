package model

import (
	"encoding/json"
	"scheduler/params/req"
	"time"
)

// Task 表示任务的数据库实体
type Task struct {
	Id         int64     `gorm:"column:id;not null;type:bigint(20);autoIncrement;primaryKey;"`
	Status     int32     `gorm:"column:status;not null;type:int(11);default:2;"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;type:datetime;default:CURRENT_TIMESTAMP;"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;type:datetime;default:CURRENT_TIMESTAMP;"`
	Name       string    `gorm:"column:name;not null;type:varchar(80);default:'';"`
	Spec       string    `gorm:"column:spec;index:idx_spec;not null;type:varchar(40);default:'';"`
	Tag        string    `gorm:"column:tag;index:idx_tag;not null;type:varchar(40);default:'';"`
	RetryMax   int32     `gorm:"column:retry_max;not null;type:int(11);default:0;"`
	RetryCycle int32     `gorm:"column:retry_cycle;not null;type:int(11);default:0;"`
	Url        string    `gorm:"column:url;type:text;"`
	BackupUrl  string    `gorm:"column:backup_url;type:text;"`
	Method     string    `gorm:"column:method;not null;type:varchar(6);default:'GET';"`
	Body       string    `gorm:"column:body;type:text;"`
	Header     string    `gorm:"column:header;type:text;"`
	Total      int64     `gorm:"column:total;not null;type:bigint(20);default:0;"`
	TimeLock   int64     `gorm:"column:time_lock;not null;type:bigint(20);default:0;"`
}

// NewTask 创建 Task 实例
func NewTask() *Task {
	return &Task{}
}

// Build 根据 TaskCommand 构建 Task
func (po *Task) Build(id int64, command req.TaskCommand) *Task {
	headerJson := ""
	if len(command.Header) > 0 {
		headerBytes, err := json.Marshal(command.Header)
		if err == nil {
			headerJson = string(headerBytes)
		}
	}
	po.Id = id
	po.Status = command.Status
	po.Name = command.Name
	po.Tag = command.Tag
	po.Spec = command.Spec
	po.RetryMax = command.RetryMax
	po.RetryCycle = command.RetryCycle
	po.Url = command.Url
	po.BackupUrl = command.BackupUrl
	po.Method = command.Method
	po.Body = command.Body
	po.Header = headerJson
	return po
}
