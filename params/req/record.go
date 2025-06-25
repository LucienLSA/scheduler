package req

import (
	"fmt"
	"scheduler/pkg/tools"
	"strings"
	"time"
)

// RecordQuery 用于记录查询参数
type RecordQuery struct {
	TaskId    int64  `json:"taskId"`
	Code      int    `json:"code"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Shard     string `json:"shard"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}

// ConversionAndVerify 校验并转换记录查询参数
func (query *RecordQuery) ConversionAndVerify() error {
	query.StartTime = strings.TrimSpace(query.StartTime)
	query.EndTime = strings.TrimSpace(query.EndTime)
	query.Shard = strings.TrimSpace(query.Shard)
	if len(query.Shard) < 7 {
		dateStr := time.Now().Format("2006-01-02")
		dateArr := strings.Split(dateStr, "-")
		query.Shard = fmt.Sprintf("%s_%s", dateArr[0], dateArr[1])
	}
	return tools.VerifyPageQueryParams(query.PageIndex, query.PageSize)
}
