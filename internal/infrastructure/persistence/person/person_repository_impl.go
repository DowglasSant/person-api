package person

import (
	"fmt"

	personModel "pessoas-api/internal/domain/person/model"
	"pessoas-api/internal/domain/person/ports"

	"gorm.io/gorm"
)

// PersonRepositoryImpl implements the ports.PersonRepository interface.
// This is the adapter for PostgreSQL database persistence.
type PersonRepositoryImpl struct {
	db *gorm.DB
}

// NewPersonRepository creates a new instance of PersonRepositoryImpl.
// It returns the implementation as the PersonRepository interface.
func NewPersonRepository(db *gorm.DB) ports.PersonRepository {
	return &PersonRepositoryImpl{
		db: db,
	}
}

func (r *PersonRepositoryImpl) Save(p *personModel.Person) (int, error) {
	entity := FromDomain(p)

	result := r.db.Create(entity)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to save person: %w", result.Error)
	}

	return entity.ID, nil
}

func (r *PersonRepositoryImpl) FindAll(page, pageSize int, sortBy, sortOrder string) ([]*personModel.Person, int64, error) {
	var entities []PersonEntity
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&PersonEntity{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count persons: %w", err)
	}

	orderClause := buildOrderClause(sortBy, sortOrder)

	result := r.db.Offset(offset).Limit(pageSize).Order(orderClause).Find(&entities)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to find persons: %w", result.Error)
	}

	persons := make([]*personModel.Person, len(entities))
	for i, entity := range entities {
		persons[i] = entity.ToDomain()
	}

	return persons, total, nil
}

func buildOrderClause(sortBy, sortOrder string) string {
	validFields := map[string]string{
		"id":         "id",
		"name":       "name",
		"cpf":        "cpf",
		"email":      "email",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}

	field, exists := validFields[sortBy]
	if !exists {
		field = "id"
	}

	order := "DESC"
	if sortOrder == "asc" || sortOrder == "ASC" {
		order = "ASC"
	}

	return fmt.Sprintf("%s %s", field, order)
}

func (r *PersonRepositoryImpl) FindByCPF(cpf string) (*personModel.Person, error) {
	var entity PersonEntity

	result := r.db.Where("cpf = ?", cpf).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find person by CPF: %w", result.Error)
	}

	return entity.ToDomain(), nil
}

func (r *PersonRepositoryImpl) FindByID(id int) (*personModel.Person, error) {
	var entity PersonEntity

	result := r.db.First(&entity, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find person by ID: %w", result.Error)
	}

	return entity.ToDomain(), nil
}

func (r *PersonRepositoryImpl) Update(p *personModel.Person) error {
	entity := FromDomain(p)

	result := r.db.Model(&PersonEntity{}).Where("id = ?", entity.ID).Updates(entity)
	if result.Error != nil {
		return fmt.Errorf("failed to update person: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("person not found")
	}

	return nil
}

func (r *PersonRepositoryImpl) Delete(id int) error {
	result := r.db.Delete(&PersonEntity{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete person: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("person not found")
	}

	return nil
}
