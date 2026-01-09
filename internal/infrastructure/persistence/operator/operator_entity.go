package operator

import (
	"time"

	operator "pessoas-api/internal/domain/operator/model"
)

type OperatorEntity struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255);not null"`
	Active       bool      `gorm:"default:true;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (OperatorEntity) TableName() string {
	return "operators"
}

func (e *OperatorEntity) ToDomain() *operator.Operator {
	return &operator.Operator{
		ID:           e.ID,
		Username:     e.Username,
		Email:        e.Email,
		PasswordHash: e.PasswordHash,
		Active:       e.Active,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func FromDomain(op *operator.Operator) *OperatorEntity {
	return &OperatorEntity{
		ID:           op.ID,
		Username:     op.Username,
		Email:        op.Email,
		PasswordHash: op.PasswordHash,
		Active:       op.Active,
		CreatedAt:    op.CreatedAt,
		UpdatedAt:    op.UpdatedAt,
	}
}
