package models

type UserTag struct {
	UserID string `gorm:"primaryKey"`
	Tag    int    `gorm:"primaryKey"`
}
