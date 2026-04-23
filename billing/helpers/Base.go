package helpers

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null;type:timestamp;autoCreateTime;column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null;type:timestamp;autoUpdateTime;column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (b *BaseModel) Get() (interface{}, error) {
	return nil, nil
}
