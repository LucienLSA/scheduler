package model

// Metadata 用于存储元数据信息
type Metadata struct {
	Id              int64 `gorm:"column:id;not null;type:bigint(20);autoIncrement;primaryKey;"`
	TaskEditVersion int64 `gorm:"column:task_edit_version;not null;type:bigint(20);default:0;"`
}
