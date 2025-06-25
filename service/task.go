package service

import (
	"context"
	"scheduler/db/model"
	"scheduler/db/mysql"
	"scheduler/params/req"
	"sync"

	"go.uber.org/zap"
)

// 单例模式
var taskSrvIns *TaskSrv
var taskSrvOnce sync.Once

type TaskSrv struct {
}

// 单例实例 不对外暴露 通过GetUserSrv来返回实例对象
func GetTaskSrv() *TaskSrv {
	taskSrvOnce.Do(func() {
		taskSrvIns = &TaskSrv{}
	})
	return taskSrvIns
}

// 重置单例 便于测试
func ResetTaskSrv() {
	taskSrvOnce = sync.Once{}
	taskSrvIns = nil
}

func (s *TaskSrv) ListTask(ctx context.Context, req *req.TaskQuery) ([]model.Task, int64, error) {
	taskDao := mysql.NewTaskDao(ctx)
	list, total, err := taskDao.ListTask(req.Name, req.Tag, req.Spec, req.Status, req.PageIndex, req.PageSize)
	if err != nil {
		zap.L().Error("taskDao.ListTask failed", zap.Error(err))
		return nil, 0, err
	}
	return list, total, nil
}
