package resp

import (
	"encoding/json"
	"scheduler/db/model"
)

// TaskDTO 用于返回任务数据到前端
type TaskDTO struct {
	Id         int64             `json:"id"`
	Status     int32             `json:"status"`
	Name       string            `json:"name"`
	Tag        string            `json:"tag"`
	Spec       string            `json:"spec"`
	RetryMax   int32             `json:"retryMax"`
	RetryCycle int32             `json:"retryCycle"`
	Url        string            `json:"url"`
	BackupUrl  string            `json:"backupUrl"`
	Method     string            `json:"method"`
	Body       string            `json:"body"`
	Header     map[string]string `json:"header"`
	Total      int64             `json:"total"`
	CreatedAt  int64             `json:"createdAt"`
	UpdatedAt  int64             `json:"updatedAt"`
}

// NewTaskDTO 创建 TaskDTO 实例
func NewTaskDTO() *TaskDTO {
	return &TaskDTO{}
}

// Build 根据 Task 构建 TaskDTO
func (dto *TaskDTO) Build(task model.Task) *TaskDTO {
	headerObj := map[string]string{}
	if len(task.Header) > 0 {
		_ = json.Unmarshal([]byte(task.Header), &headerObj)
	}
	dto.Id = task.Id
	dto.Status = task.Status
	dto.Name = task.Name
	dto.Tag = task.Tag
	dto.Spec = task.Spec
	dto.RetryMax = task.RetryMax
	dto.RetryCycle = task.RetryCycle
	dto.Url = task.Url
	dto.BackupUrl = task.BackupUrl
	dto.Method = task.Method
	dto.Body = task.Body
	dto.Header = headerObj
	dto.Total = task.Total
	dto.CreatedAt = task.CreatedAt.UnixMilli()
	dto.UpdatedAt = task.UpdatedAt.UnixMilli()
	return dto
}

// BuildWithTask 根据 Task 列表构建分页 DTO
func (dto *PageDTO) BuildWithTask(list []model.Task, total int64) *PageDTO {
	var dtoList []TaskDTO
	if len(list) > 0 {
		for _, v := range list {
			dtoItem := TaskDTO{}
			dtoList = append(dtoList, *dtoItem.Build(v))
		}
	}
	dto.Total = total
	dto.List = dtoList
	return dto
}
