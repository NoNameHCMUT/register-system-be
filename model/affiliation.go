package model

type Affiliation struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	StdName string `gorm:"uniqueIndex;not null" json:"std_name"`
}
