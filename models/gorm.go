package models

import "time"

type Base struct {
	ID        *int64     `json:"id" gorm:"column:id;primary_key" validate:"-"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at" validate:"-"`
	CreatedBy *int64     `json:"created_by" gorm:"column:created_by" validate:"-"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at" validate:"-"`
	UpdatedBy *int64     `json:"updated_by" gorm:"column:updated_by" validate:"-"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at" validate:"-"`
}
