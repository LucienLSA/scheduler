package req

import (
	"errors"
	"scheduler/pkg/tools"
	"strings"
)

// TaskCommand 用于接收任务相关的请求参数
type TaskCommand struct {
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
}

// ConversionAndVerifyWithAdd 校验并转换新增任务参数
func (cmd *TaskCommand) ConversionAndVerifyWithAdd() error {
	cmd.Spec = strings.TrimSpace(cmd.Spec)
	cmd.Tag = strings.TrimSpace(cmd.Tag)
	cmd.Name = strings.TrimSpace(cmd.Name)
	cmd.Url = strings.TrimSpace(cmd.Url)
	if len(cmd.Tag) == 0 {
		cmd.Tag = "default"
	}
	if len(cmd.Name) == 0 {
		return errors.New("name is empty")
	}
	if len(cmd.Url) == 0 {
		return errors.New("url is empty")
	}
	if !tools.VerifyUrl(cmd.Url) {
		return errors.New("url format is incorrect")
	}
	if len(cmd.BackupUrl) > 0 && !tools.VerifyUrl(cmd.BackupUrl) {
		return errors.New("backup url format is incorrect")
	}
	if len(cmd.Spec) == 0 {
		return errors.New("spec is empty")
	}
	if cmd.Method != "GET" && cmd.Method != "POST" && cmd.Method != "PUT" && cmd.Method != "PATCH" && cmd.Method != "DELETE" {
		return errors.New("method is not match")
	}
	return nil
}

// ConversionAndVerifyWithEdit 校验并转换编辑任务参数
func (cmd *TaskCommand) ConversionAndVerifyWithEdit() error {
	cmd.Spec = strings.TrimSpace(cmd.Spec)
	cmd.Tag = strings.TrimSpace(cmd.Tag)
	cmd.Name = strings.TrimSpace(cmd.Name)
	cmd.Url = strings.TrimSpace(cmd.Url)
	if len(cmd.Url) > 0 && !tools.VerifyUrl(cmd.Url) {
		return errors.New("url format is incorrect")
	}
	if len(cmd.BackupUrl) > 0 && cmd.BackupUrl != "nil" && !tools.VerifyUrl(cmd.BackupUrl) {
		return errors.New("backup url format is incorrect")
	}
	if len(cmd.Method) > 0 && cmd.Method != "GET" && cmd.Method != "POST" && cmd.Method != "PUT" && cmd.Method != "PATCH" && cmd.Method != "DELETE" {
		return errors.New("method is not match")
	}
	return nil
}

// TaskQuery 用于任务查询参数
type TaskQuery struct {
	Status    int    `json:"status"`
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	Spec      string `json:"spec"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}

// ConversionAndVerify 校验并转换任务查询参数
func (query *TaskQuery) ConversionAndVerify() error {
	query.Name = strings.TrimSpace(query.Name)
	query.Tag = strings.TrimSpace(query.Tag)
	query.Spec = strings.TrimSpace(query.Spec)
	return tools.VerifyPageQueryParams(query.PageIndex, query.PageSize)
}
