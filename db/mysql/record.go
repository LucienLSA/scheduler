package mysql

import (
	"context"
	"fmt"
	"scheduler/db/model"

	"gorm.io/gorm"
)

type RecordDao struct {
	*gorm.DB
}

func NewRecordDao(ctx context.Context) *RecordDao {
	return &RecordDao{NewDBClient(ctx)}
}

func NewRecordDaoByDB(db *gorm.DB) *RecordDao {
	return &RecordDao{db}
}
func (dao *RecordDao) AddRecord(record model.Record) error {
	err := dao.DB.AutoMigrate(&model.Record{})
	if err != nil {
		return err
	}
	return dao.DB.Model(&model.Record{}).Create(&record).Error
}

func (dao *RecordDao) ListRecord(shard string, taskId int64, code int, startTime string, endTime string, pageIndex int, pageSize int) ([]model.Record, int64, error) {
	var recordList []model.Record
	var total int64
	sql := dao.DB.Table(fmt.Sprintf("record_%s", shard)).Where("task_id = ?", taskId)
	if code != 0 {
		sql.Where("code = ?", code)
	}
	if len(startTime) > 0 {
		sql.Where("executed_at >= ?", startTime)
	}
	if len(endTime) > 0 {
		sql.Where("executed_at <= ?", endTime)
	}
	sql.Count(&total)
	sql = sql.Order("id desc").Limit(pageSize).Offset((pageIndex - 1) * pageSize)
	err := sql.Find(&recordList).Error
	return recordList, total, err
}
