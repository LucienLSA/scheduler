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

func (s *TaskSrv) ListStartedTaskBySpec(ctx context.Context, spec string) ([]model.Task, error) {
	taskDao := mysql.NewTaskDao(ctx)
	return taskDao.ListStartedTaskBySpec(spec)
}

func (s *TaskSrv) TryExecuteTask(ctx context.Context, task model.Task) (int64, error) {
	taskDao := mysql.NewTaskDao(ctx)
	return taskDao.TryExecuteTask(task)
}

func (s *TaskSrv) AddTask(ctx context.Context, task model.Task) (int64, error) {
	taskDao := mysql.NewTaskDao(ctx)
	metadataDao := mysql.NewMetadataDao(ctx)
	err := metadataDao.ChangeTaskEditVersion()
	if err != nil {
		return 0, err
	}
	return taskDao.AddTask(task)
}

func (s *TaskSrv) GetTask(ctx context.Context, id int64) (model.Task, error) {
	taskDao := mysql.NewTaskDao(ctx)
	return taskDao.GetTask(id)
}

func (s *TaskSrv) EditTask(ctx context.Context, task model.Task) error {
	taskDao := mysql.NewTaskDao(ctx)
	metadataDao := mysql.NewMetadataDao(ctx)
	err := metadataDao.ChangeTaskEditVersion()
	if err != nil {
		return err
	}
	return taskDao.EditTask(task)
}

func (s *TaskSrv) DeleteTask(ctx context.Context, id int64) error {
	taskDao := mysql.NewTaskDao(ctx)
	return taskDao.DeleteTask(id)
}
