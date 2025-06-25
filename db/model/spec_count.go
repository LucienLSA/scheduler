package model

// SpecCount 用于统计 Spec 数量
type SpecCount struct {
	Spec   string `gorm:"column:spec;" json:"spec"`
	Number int64  `gorm:"column:number;" json:"number"`
}
