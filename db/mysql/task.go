package mysql

import (
	"context"
	"scheduler/db/model"

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
