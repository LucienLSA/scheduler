package service

import (
	"context"
	"scheduler/db/mysql"
	"sync"
)

// 单例模式
var metadataSrvIns *MetadataSrv
var metadataSrvOnce sync.Once

type MetadataSrv struct {
}

// 单例实例 不对外暴露 通过GetUserSrv来返回实例对象
func GetMetadataSrv() *MetadataSrv {
	metadataSrvOnce.Do(func() {
		metadataSrvIns = &MetadataSrv{}
	})
	return metadataSrvIns
}

// 重置单例 便于测试
func ResetMetadataSrv() {
	metadataSrvOnce = sync.Once{}
	metadataSrvIns = nil
}

func (s *MetadataSrv) GetTaskEditVersion(ctx context.Context) (int64, error) {
	metadataDao := mysql.NewMetadataDao(ctx)
	return metadataDao.GetTaskEditVersion()
}
