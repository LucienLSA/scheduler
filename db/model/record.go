package model

import (
	"fmt"
	"strings"
	"time"
)

// Record 表示任务执行记录的数据库实体
type Record struct {
	Id         int64     `gorm:"column:id;not null;type:bigint(20);autoIncrement;primaryKey;"`
	TaskId     int64     `gorm:"column:task_id;index:idx_task_id;not null;type:bigint(20);default:0;"`
	ExecutedAt time.Time `gorm:"column:executed_at;index:idx_executed_at;not null;type:datetime;default:CURRENT_TIMESTAMP;"`
	Result     string    `gorm:"column:result;type:text;"`
	Code       int32     `gorm:"column:code;not null;type:int(11);default:0;"`
	TimeCost   int32     `gorm:"column:time_cost;not null;type:int(11);default:0;"`
	RetryCount int32     `gorm:"column:retry_count;not null;type:int(11);default:0;"`
	IsBackup   int32     `gorm:"column:is_backup;not null;type:int(11);default:0;"`
}

// TableName 返回分表后的记录表名
func (Record) TableName() string {
	dateStr := time.Now().Format("2006-01-02")
	dateArr := strings.Split(dateStr, "-")
	return fmt.Sprintf("record_%s_%s", dateArr[0], dateArr[1])
}
