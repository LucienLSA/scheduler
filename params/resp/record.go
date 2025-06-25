package resp

// RecordDTO 用于返回记录数据到前端
type RecordDTO struct {
	Id         int64  `json:"id"`
	TaskId     int64  `json:"taskId"`
	ExecutedAt int64  `json:"executedAt"`
	Result     string `json:"result"`
	TimeCost   int32  `json:"timeCost"`
	Code       int32  `json:"code"`
	IsBackup   int32  `json:"isBackup"`
	RetryCount int32  `json:"retryCount"`
}
