package resp

import "scheduler/db/model"

// PageDTO 用于分页返回数据的结构体
// Total 表示总数，List 表示数据列表
type PageDTO struct {
	Total int64 `json:"total"`
	List  any   `json:"list"`
}

func NewPageDTO() *PageDTO {
	return &PageDTO{}
}

// BuildWithRecord 根据 Record 列表构建分页 DTO
func (dto *PageDTO) BuildWithRecord(list []model.Record, total int64) *PageDTO {
	var dtoList []RecordDTO
	if len(list) > 0 {
		for _, v := range list {
			dtoList = append(dtoList, RecordDTO{
				Id:         v.Id,
				TaskId:     v.TaskId,
				Result:     v.Result,
				Code:       v.Code,
				IsBackup:   v.IsBackup,
				TimeCost:   v.TimeCost,
				RetryCount: v.RetryCount,
				ExecutedAt: v.ExecutedAt.UnixMilli(),
			})
		}
	}
	dto.Total = total
	dto.List = dtoList
	return dto
}
