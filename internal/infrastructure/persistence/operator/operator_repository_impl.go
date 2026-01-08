package operator

import (
	"errors"
	"log"

	operator "pessoas-api/internal/domain/operator/model"
	"pessoas-api/internal/domain/operator/ports"

	"gorm.io/gorm"
)

type OperatorRepositoryImpl struct {
	db *gorm.DB
}

func NewOperatorRepository(db *gorm.DB) ports.OperatorRepository {
	return &OperatorRepositoryImpl{db: db}
}

func (r *OperatorRepositoryImpl) Save(op *operator.Operator) (int, error) {
	entity := FromDomain(op)

	result := r.db.Create(entity)
	if result.Error != nil {
		log.Printf("[ERROR] OperatorRepository.Save - Failed to save operator: %v", result.Error)
		return 0, result.Error
	}

	return entity.ID, nil
}

func (r *OperatorRepositoryImpl) FindByUsername(username string) (*operator.Operator, error) {
	var entity OperatorEntity

	result := r.db.Where("username = ?", username).First(&entity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("[ERROR] OperatorRepository.FindByUsername - Failed to find operator: %v", result.Error)
		return nil, result.Error
	}

	return entity.ToDomain(), nil
}

func (r *OperatorRepositoryImpl) FindByEmail(email string) (*operator.Operator, error) {
	var entity OperatorEntity

	result := r.db.Where("email = ?", email).First(&entity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("[ERROR] OperatorRepository.FindByEmail - Failed to find operator: %v", result.Error)
		return nil, result.Error
	}

	return entity.ToDomain(), nil
}

func (r *OperatorRepositoryImpl) FindByID(id int) (*operator.Operator, error) {
	var entity OperatorEntity

	result := r.db.First(&entity, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("[ERROR] OperatorRepository.FindByID - Failed to find operator: %v", result.Error)
		return nil, result.Error
	}

	return entity.ToDomain(), nil
}
