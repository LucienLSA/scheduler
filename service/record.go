package service

import (
	"context"
	"scheduler/db/model"
	"scheduler/db/mysql"
	"sync"
)

// 单例模式
var recordSrvIns *RecordSrv
var recordSrvOnce sync.Once

type RecordSrv struct {
}

// 单例实例 不对外暴露 通过GetUserSrv来返回实例对象
func GetRecordSrv() *RecordSrv {
	recordSrvOnce.Do(func() {
		recordSrvIns = &RecordSrv{}
	})
	return recordSrvIns
}

// 重置单例 便于测试
func ResetRecordSrv() {
	recordSrvOnce = sync.Once{}
	recordSrvIns = nil
}

func (s *RecordSrv) AddRecord(ctx context.Context, record model.Record) error {
	recordDao := mysql.NewRecordDao(ctx)
	return recordDao.AddRecord(record)
}

func (s *RecordSrv) ListRecord(ctx context.Context, shard string, taskId int64,
	code int, startTime string, endTime string, pageIndex int, pageSize int) ([]model.Record, int64, error) {
	recordDao := mysql.NewRecordDao(ctx)
	return recordDao.ListRecord(shard, taskId, code, startTime,
		endTime, pageIndex, pageSize)
}
