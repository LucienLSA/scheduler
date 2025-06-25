package mysql

import (
	"context"
	"scheduler/db/model"

	"gorm.io/gorm"
)

type MetadataDao struct {
	*gorm.DB
}

func NewMetadataDao(ctx context.Context) *MetadataDao {
	return &MetadataDao{NewDBClient(ctx)}
}

func NewMetadataDaoByDB(db *gorm.DB) *MetadataDao {
	return &MetadataDao{db}
}
func (dao *MetadataDao) GetTaskEditVersion() (int64, error) {
	var metadata model.Metadata
	err := dao.DB.Model(&model.Metadata{}).Where("id = ?", 1).First(&metadata).Error
	if err != nil {
		return 0, err
	}
	return metadata.TaskEditVersion, nil
}

func (dao *MetadataDao) ChangeTaskEditVersion() error {
	return dao.DB.Table("metadata").Where("id = ?", 1).UpdateColumn("task_edit_version", gorm.Expr("task_edit_version + ?", 1)).Error
}
