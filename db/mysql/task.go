package mysql

import (
	"context"
	"scheduler/config"
	"scheduler/db/model"
	"time"

	"gorm.io/gorm"
)

type TaskDao struct {
	*gorm.DB
}

func NewTaskDao(ctx context.Context) *TaskDao {
	return &TaskDao{NewDBClient(ctx)}
}

func NewTaskDaoByDB(db *gorm.DB) *TaskDao {
	return &TaskDao{db}
}

func (dao *TaskDao) ListTask(name string, tag string, spec string, status int, pageIndex int, pageSize int) ([]model.Task, int64, error) {
	var taskList []model.Task
	var total int64

	// 构建基础查询
	baseQuery := dao.DB.Model(&model.Task{})
	if status > 0 {
		baseQuery = baseQuery.Where("status = ?", status)
	}
	if len(name) > 0 {
		baseQuery = baseQuery.Where("name like ?", "%"+name+"%")
	}
	if len(tag) > 0 {
		baseQuery = baseQuery.Where("tag = ?", tag)
	}
	if len(spec) > 0 {
		baseQuery = baseQuery.Where("spec = ?", spec)
	}

	// 先查询总数
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再查询分页数据
	err := baseQuery.Order("id desc").Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&taskList).Error
	return taskList, total, err
}

func (dao *TaskDao) ListStartedTaskBySpec(spec string) ([]model.Task, error) {
	var taskList []model.Task
	err := dao.DB.Model(&model.Task{}).Where("spec = ? and status = ?", spec, 1).Find(&taskList).Error
	return taskList, err
}

func (dao *TaskDao) TryExecuteTask(task model.Task) (int64, error) {
	sql := dao.DB.Table("task").Where("id = ? and total = ?", task.Id, task.Total).Updates(map[string]interface{}{
		"total":     gorm.Expr("total + ?", 1),
		"time_lock": time.Now().UnixMilli() + config.Conf.DBConfig.ExecutedLockTime,
	})
	if sql.Error != nil {
		return 0, sql.Error
	}
	return sql.RowsAffected, nil
}

func (dao *TaskDao) AddTask(task model.Task) (int64, error) {
	task.UpdatedAt = time.Now()
	task.CreatedAt = task.UpdatedAt
	return task.Id, dao.DB.Table("task").Create(&task).Error
}

func (dao *TaskDao) GetTask(id int64) (model.Task, error) {
	var task model.Task
	err := dao.DB.Model(&model.Task{}).Where("id = ?", id).First(&task).Error
	if err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (dao *TaskDao) EditTask(task model.Task) error {
	task.UpdatedAt = time.Now()
	if task.BackupUrl == "nil" {
		err := dao.DB.Table("task").Where("id = ?", task.Id).Updates(map[string]any{
			"backup_url": "",
			"updated_at": task.UpdatedAt,
		}).Error
		if err != nil {
			return err
		}
		task.BackupUrl = ""
	}
	return dao.DB.Table("task").Where("id = ?", task.Id).Updates(&task).Error
}

func (dao *TaskDao) DeleteTask(id int64) error {
	return dao.DB.Table("task").Delete(model.Task{}, id).Error
}
