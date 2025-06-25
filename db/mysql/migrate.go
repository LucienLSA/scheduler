package mysql

import "scheduler/db/model"

func migrate() (err error) {
	err = _db.AutoMigrate(&model.Task{}, &model.Metadata{},
		&model.SpecCount{}, &model.Record{}, &model.TagCount{})
	if err != nil {
		return
	}
	return nil
}
