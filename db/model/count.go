package model

// TagCount 用于统计 Tag 数量
type TagCount struct {
	Tag    string `gorm:"column:tag;" json:"tag"`
	Number int64  `gorm:"column:number;" json:"number"`
}

// SpecCount 用于统计 Spec 数量
type SpecCount struct {
	Spec   string `gorm:"column:spec;" json:"spec"`
	Number int64  `gorm:"column:number;" json:"number"`
}
