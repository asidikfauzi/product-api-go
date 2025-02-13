package model

import (
	"github.com/google/uuid"
	"time"
)

type Categories struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name      string     `gorm:"unique;type:varchar(100);not null" json:"name"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy uuid.UUID  `gorm:"type:uuid" json:"updated_by"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
	DeletedBy *uuid.UUID `gorm:"type:uuid" json:"deleted_by"`
	// Relation
	Products []Products `gorm:"many2many:category_has_products" json:"products"`
}
