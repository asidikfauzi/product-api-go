package model

import "github.com/google/uuid"

type ProductMeasurements struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name string    `gorm:"type:varchar(50);not null" json:"name"`
	// Relation
	Product []Products `gorm:"foreignKey:ProductMeasurementID;references:ID"`
}
