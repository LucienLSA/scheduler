package mysql

import (
	"scheduler/db/model"

	"go.uber.org/zap"
)

func migrate() (err error) {
	err = _db.AutoMigrate(&model.Task{}, &model.Metadata{},
		&model.SpecCount{}, &model.Record{}, &model.TagCount{})
	if err != nil {
		zap.L().Error("DB AutoMigrate failed, err:%v\n", zap.Error(err))
		return
	}
	return nil
}
