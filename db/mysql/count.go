package mysql

import (
	"context"
	"scheduler/db/model"

	"gorm.io/gorm"
)

type CountDao struct {
	*gorm.DB
}

func NewCountDao(ctx context.Context) *CountDao {
	return &CountDao{NewDBClient(ctx)}
}

func NewCountDaoByDB(db *gorm.DB) *CountDao {
	return &CountDao{db}
}

func (dao *CountDao) ListTagCount(status int) ([]model.TagCount, error) {
	var tagCountList []model.TagCount
	sql := dao.DB.Table("task").Select("tag", "count(*) as total")
	if status > 0 {
		sql = sql.Where("status = ?", status)
	}
	err := sql.Group("tag").Find(&tagCountList).Error
	if err != nil {
		return nil, err
	}
	return tagCountList, err
}

func (dao *CountDao) ListSpecCount(status int) ([]model.SpecCount, error) {
	var specCountList []model.SpecCount
	sql := dao.DB.Table("task").Select("spec", "count(*) as number")
	if status > 0 {
		sql = sql.Where("status = ?", status)
	}
	err := sql.Group("spec").Find(&specCountList).Error
	if err != nil {
		return nil, err
	}
	return specCountList, err
}

func (dao *CountDao) ListStartedSpec() ([]string, error) {
	var specList []string
	specCountList, err := dao.ListSpecCount(1)
	if err != nil {
		return nil, err
	}
	for _, v := range specCountList {
		specList = append(specList, v.Spec)
	}
	return specList, err
}
