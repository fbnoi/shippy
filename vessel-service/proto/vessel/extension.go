package vessel

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (model *Vessel) BeforeCreate(tx *gorm.DB) (error) {
	model.Id = uuid.New().String()
	return nil
}