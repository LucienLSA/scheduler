package service

import (
	"context"
	"scheduler/db/model"
	"scheduler/db/mysql"
	"sync"
)

// 单例模式
var countSrvIns *CountSrv
var countSrvOnce sync.Once

type CountSrv struct {
}

// 单例实例 不对外暴露 通过GetUserSrv来返回实例对象
func GetCountSrv() *CountSrv {
	countSrvOnce.Do(func() {
		countSrvIns = &CountSrv{}
	})
	return countSrvIns
}

// 重置单例 便于测试
func ResetCountSrv() {
	countSrvOnce = sync.Once{}
	countSrvIns = nil
}

func (s *CountSrv) ListStartedSpec(ctx context.Context) ([]string, error) {
	countDao := mysql.NewCountDao(ctx)
	return countDao.ListStartedSpec()
}

func (s *CountSrv) ListTagCount(ctx context.Context, status int) ([]model.TagCount, error) {
	countDao := mysql.NewCountDao(ctx)
	return countDao.ListTagCount(status)
}

func (s *CountSrv) ListSpecCount(ctx context.Context, status int) ([]model.SpecCount, error) {
	countDao := mysql.NewCountDao(ctx)
	return countDao.ListSpecCount(status)
}
