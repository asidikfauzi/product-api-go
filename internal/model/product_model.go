package model

import (
	"github.com/google/uuid"
	"time"
)

type Products struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProductMeasurementID uuid.UUID  `gorm:"type:uuid;index" json:"product_measurement_id"`
	Name                 string     `gorm:"varchar(100);not null" json:"name"`
	Description          *string    `gorm:"text" json:"description"`
	PurchasePrice        int        `gorm:"not null" json:"purchase_price"`
	SellingPrice         int        `gorm:"not null" json:"selling_price"`
	TotalStock           int        `gorm:"not null" json:"total_stock"`
	MinimumStock         int        `gorm:"default:10;not null" json:"minimum_stock"`
	Image                *string    `gorm:"varchar(255)" json:"image"`
	CreatedAt            time.Time  `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy            uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy            uuid.UUID  `gorm:"type:uuid" json:"updated_by"`
	DeletedAt            *time.Time `gorm:"index" json:"deleted_at"`
	DeletedBy            *uuid.UUID `gorm:"type:uuid" json:"deleted_by"`
	// Relation
	ProductMeasurement ProductMeasurements `gorm:"foreignKey:ProductMeasurementID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"product_measurement"`
	Categories         []Categories        `gorm:"many2many:category_has_products" json:"categories"`
}
