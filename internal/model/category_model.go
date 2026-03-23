package model

import "github.com/google/uuid"

type CategoryModel struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
}

func (CategoryModel) TableName() string {
	return "categories"
}
