package person

import (
	"fmt"
	"testing"
	"time"

	personModel "pessoas-api/internal/domain/person/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type testPersonEntity struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name;type:varchar(255);not null"`
	CPF         string    `gorm:"column:cpf;type:varchar(11);not null;uniqueIndex"`
	BirthDate   time.Time `gorm:"column:birth_date;type:date;not null"`
	PhoneNumber string    `gorm:"column:phone_number;type:varchar(11);not null"`
	Email       string    `gorm:"column:email;type:varchar(255);not null"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;not null"`
}

func (testPersonEntity) TableName() string {
	return "person"
}

type testPersonRepositoryImpl struct {
	db *gorm.DB
}

func newTestPersonRepository(db *gorm.DB) *testPersonRepositoryImpl {
	return &testPersonRepositoryImpl{db: db}
}

func (r *testPersonRepositoryImpl) Save(p *personModel.Person) (int, error) {
	entity := &testPersonEntity{
		Name:        p.Name,
		CPF:         p.CPF,
		BirthDate:   p.BirthDate,
		PhoneNumber: p.PhoneNumber,
		Email:       p.Email,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	result := r.db.Create(entity)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to save person: %w", result.Error)
	}

	return entity.ID, nil
}

func (r *testPersonRepositoryImpl) FindAll(page, pageSize int, sortBy, sortOrder string) ([]*personModel.Person, int64, error) {
	var entities []testPersonEntity
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&testPersonEntity{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count persons: %w", err)
	}

	orderClause := buildOrderClause(sortBy, sortOrder)

	result := r.db.Offset(offset).Limit(pageSize).Order(orderClause).Find(&entities)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to find persons: %w", result.Error)
	}

	persons := make([]*personModel.Person, len(entities))
	for i, entity := range entities {
		persons[i] = &personModel.Person{
			ID:          entity.ID,
			Name:        entity.Name,
			CPF:         entity.CPF,
			BirthDate:   entity.BirthDate,
			PhoneNumber: entity.PhoneNumber,
			Email:       entity.Email,
			CreatedAt:   entity.CreatedAt,
			UpdatedAt:   entity.UpdatedAt,
		}
	}

	return persons, total, nil
}

func (r *testPersonRepositoryImpl) FindByCPF(cpf string) (*personModel.Person, error) {
	var entity testPersonEntity

	result := r.db.Where("cpf = ?", cpf).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find person by CPF: %w", result.Error)
	}

	return &personModel.Person{
		ID:          entity.ID,
		Name:        entity.Name,
		CPF:         entity.CPF,
		BirthDate:   entity.BirthDate,
		PhoneNumber: entity.PhoneNumber,
		Email:       entity.Email,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}, nil
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&testPersonEntity{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func createValidPerson(t *testing.T) *personModel.Person {
	person, err := personModel.NewPerson(
		"John Doe",
		"11144477735",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"john.doe@example.com",
	)
	if err != nil {
		t.Fatalf("failed to create valid person: %v", err)
	}
	return person
}

func TestPersonRepositoryImpl_Save_Success(t *testing.T) {
	assert := assert.New(t)
	db := setupTestDB(t)

	repo := newTestPersonRepository(db)
	person := createValidPerson(t)

	id, err := repo.Save(person)

	assert.NoError(err)
	assert.NotZero(id)
	assert.Greater(id, 0)

	var savedEntity testPersonEntity
	result := db.First(&savedEntity, id)

	assert.NoError(result.Error)
	assert.Equal(id, savedEntity.ID)
	assert.Equal("John Doe", savedEntity.Name)
	assert.Equal("11144477735", savedEntity.CPF)
	assert.Equal("81912345678", savedEntity.PhoneNumber)
	assert.Equal("john.doe@example.com", savedEntity.Email)
	assert.Equal(time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), savedEntity.BirthDate)
}

func TestPersonRepositoryImpl_Save_MultiplePersons(t *testing.T) {
	assert := assert.New(t)
	db := setupTestDB(t)

	repo := newTestPersonRepository(db)

	person1, _ := personModel.NewPerson(
		"Alice Smith",
		"11144477735",
		time.Date(1985, time.March, 15, 0, 0, 0, 0, time.UTC),
		"81987654321",
		"alice@example.com",
	)

	person2, _ := personModel.NewPerson(
		"Bob Johnson",
		"22233344405",
		time.Date(1992, time.July, 22, 0, 0, 0, 0, time.UTC),
		"81998765432",
		"bob@example.com",
	)

	id1, err1 := repo.Save(person1)
	id2, err2 := repo.Save(person2)

	assert.NoError(err1)
	assert.NoError(err2)
	assert.NotEqual(id1, id2)
	assert.Greater(id1, 0)
	assert.Greater(id2, 0)

	var count int64
	db.Model(&testPersonEntity{}).Count(&count)
	assert.Equal(int64(2), count)
}

func TestPersonRepositoryImpl_Save_DuplicateCPF_ShouldFail(t *testing.T) {
	assert := assert.New(t)
	db := setupTestDB(t)

	repo := newTestPersonRepository(db)

	person1, _ := personModel.NewPerson(
		"Alice Smith",
		"11144477735",
		time.Date(1985, time.March, 15, 0, 0, 0, 0, time.UTC),
		"81987654321",
		"alice@example.com",
	)

	person2, _ := personModel.NewPerson(
		"Bob Johnson",
		"11144477735",
		time.Date(1992, time.July, 22, 0, 0, 0, 0, time.UTC),
		"81998765432",
		"bob@example.com",
	)

	id1, err1 := repo.Save(person1)
	assert.NoError(err1)
	assert.Greater(id1, 0)

	id2, err2 := repo.Save(person2)
	assert.Error(err2)
	assert.Zero(id2)
	assert.Contains(err2.Error(), "failed to save person")
}

func TestPersonRepositoryImpl_Save_PreservesTimestamps(t *testing.T) {
	assert := assert.New(t)
	db := setupTestDB(t)

	repo := newTestPersonRepository(db)

	createdAt := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, time.January, 2, 12, 0, 0, 0, time.UTC)

	person, _ := personModel.NewPerson(
		"Jane Doe",
		"11144477735",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"jane@example.com",
	)

	person.CreatedAt = createdAt
	person.UpdatedAt = updatedAt

	id, err := repo.Save(person)

	assert.NoError(err)
	assert.Greater(id, 0)

	var savedEntity testPersonEntity
	db.First(&savedEntity, id)

	assert.True(savedEntity.CreatedAt.Equal(createdAt))
	assert.True(savedEntity.UpdatedAt.Equal(updatedAt))
}

func TestFromDomain_ConvertsCorrectly(t *testing.T) {
	assert := assert.New(t)

	person := createValidPerson(t)
	entity := FromDomain(person)

	assert.NotNil(entity)
	assert.Equal(person.Name, entity.Name)
	assert.Equal(person.CPF, entity.CPF)
	assert.Equal(person.BirthDate, entity.BirthDate)
	assert.Equal(person.PhoneNumber, entity.PhoneNumber)
	assert.Equal(person.Email, entity.Email)
	assert.Equal(person.CreatedAt, entity.CreatedAt)
	assert.Equal(person.UpdatedAt, entity.UpdatedAt)
	assert.Zero(entity.ID)
}

func TestPersonEntity_ToDomain_ConvertsCorrectly(t *testing.T) {
	assert := assert.New(t)

	entity := &PersonEntity{
		ID:          1,
		Name:        "John Doe",
		CPF:         "11144477735",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "81912345678",
		Email:       "john@example.com",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	person := entity.ToDomain()

	assert.NotNil(person)
	assert.Equal(entity.Name, person.Name)
	assert.Equal(entity.CPF, person.CPF)
	assert.Equal(entity.BirthDate, person.BirthDate)
	assert.Equal(entity.PhoneNumber, person.PhoneNumber)
	assert.Equal(entity.Email, person.Email)
	assert.Equal(entity.CreatedAt, person.CreatedAt)
	assert.Equal(entity.UpdatedAt, person.UpdatedAt)
}

func TestPersonEntity_TableName(t *testing.T) {
	assert := assert.New(t)

	entity := PersonEntity{}
	tableName := entity.TableName()

	assert.Equal("people.person", tableName)
}

func TestPersonRepositoryImpl_FindAll_Pagination(t *testing.T) {
	assert := assert.New(t)
	db := setupTestDB(t)

	repo := newTestPersonRepository(db)

	// Create 3 persons with different CPFs
	person1, _ := personModel.NewPerson(
		"Person 1",
		"11144477735",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"person1@example.com",
	)
	person2, _ := personModel.NewPerson(
		"Person 2",
		"22233344405",
		time.Date(1990, time.January, 2, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"person2@example.com",
	)
	person3, _ := personModel.NewPerson(
		"Person 3",
		"52998224725",
		time.Date(1990, time.January, 3, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"person3@example.com",
	)

	repo.Save(person1)
	repo.Save(person2)
	repo.Save(person3)

	// Test first page with page size 2
	persons, total, err := repo.FindAll(1, 2, "id", "asc")
	assert.NoError(err)
	assert.Equal(int64(3), total)
	assert.Len(persons, 2)
	assert.Equal("Person 1", persons[0].Name)
	assert.Equal("Person 2", persons[1].Name)

	// Test second page with page size 2
	persons, total, err = repo.FindAll(2, 2, "id", "asc")
	assert.NoError(err)
	assert.Equal(int64(3), total)
	assert.Len(persons, 1)
	assert.Equal("Person 3", persons[0].Name)

	// Test all in one page
	persons, total, err = repo.FindAll(1, 10, "id", "asc")
	assert.NoError(err)
	assert.Equal(int64(3), total)
	assert.Len(persons, 3)
}

func TestPersonRepositoryImpl_FindAll_SortByName_Asc(t *testing.T) {
	assert := assert.New(t)
	db := setupTestDB(t)

	repo := newTestPersonRepository(db)

	person1, _ := personModel.NewPerson(
		"Charlie",
		"11144477735",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"charlie@example.com",
	)
	person2, _ := personModel.NewPerson(
		"Alice",
		"22233344405",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"alice@example.com",
	)
	person3, _ := personModel.NewPerson(
		"Bob",
		"52998224725",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"bob@example.com",
	)

	repo.Save(person1)
	repo.Save(person2)
	repo.Save(person3)

	persons, total, err := repo.FindAll(1, 10, "name", "asc")

	assert.NoError(err)
	assert.Equal(int64(3), total)
	assert.Len(persons, 3)
	assert.Equal("Alice", persons[0].Name)
	assert.Equal("Bob", persons[1].Name)
	assert.Equal("Charlie", persons[2].Name)
}

func TestPersonRepositoryImpl_FindAll_SortByName_Desc(t *testing.T) {
	assert := assert.New(t)
	db := setupTestDB(t)

	repo := newTestPersonRepository(db)

	person1, _ := personModel.NewPerson(
		"Charlie",
		"11144477735",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"charlie@example.com",
	)
	person2, _ := personModel.NewPerson(
		"Alice",
		"22233344405",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"alice@example.com",
	)
	person3, _ := personModel.NewPerson(
		"Bob",
		"52998224725",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"bob@example.com",
	)

	repo.Save(person1)
	repo.Save(person2)
	repo.Save(person3)

	persons, total, err := repo.FindAll(1, 10, "name", "desc")

	assert.NoError(err)
	assert.Equal(int64(3), total)
	assert.Len(persons, 3)
	assert.Equal("Charlie", persons[0].Name)
	assert.Equal("Bob", persons[1].Name)
	assert.Equal("Alice", persons[2].Name)
}

func TestPersonRepositoryImpl_FindAll_SortByInvalidField_DefaultsToID(t *testing.T) {
	assert := assert.New(t)
	db := setupTestDB(t)

	repo := newTestPersonRepository(db)

	person1, _ := personModel.NewPerson(
		"Person 1",
		"11144477735",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"person1@example.com",
	)
	person2, _ := personModel.NewPerson(
		"Person 2",
		"22233344405",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"person2@example.com",
	)

	repo.Save(person1)
	repo.Save(person2)

	// Using invalid sort field, should default to "id" with "desc"
	persons, total, err := repo.FindAll(1, 10, "invalid_field", "desc")

	assert.NoError(err)
	assert.Equal(int64(2), total)
	assert.Len(persons, 2)
	// Since default is "id desc", person2 should be first (higher ID)
	assert.Equal("Person 2", persons[0].Name)
	assert.Equal("Person 1", persons[1].Name)
}

func TestPersonRepositoryImpl_FindAll_EmptyDatabase(t *testing.T) {
	assert := assert.New(t)
	db := setupTestDB(t)

	repo := newTestPersonRepository(db)

	persons, total, err := repo.FindAll(1, 10, "id", "asc")

	assert.NoError(err)
	assert.Equal(int64(0), total)
	assert.Len(persons, 0)
}

func TestPersonRepositoryImpl_FindByCPF_Success(t *testing.T) {
	assert := assert.New(t)
	db := setupTestDB(t)

	repo := newTestPersonRepository(db)

	person, _ := personModel.NewPerson(
		"John Doe",
		"11144477735",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"81912345678",
		"john.doe@example.com",
	)

	savedID, _ := repo.Save(person)

	found, err := repo.FindByCPF("11144477735")

	assert.NoError(err)
	assert.NotNil(found)
	assert.Equal(savedID, found.ID)
	assert.Equal("John Doe", found.Name)
	assert.Equal("11144477735", found.CPF)
	assert.Equal("81912345678", found.PhoneNumber)
	assert.Equal("john.doe@example.com", found.Email)
}

func TestPersonRepositoryImpl_FindByCPF_NotFound(t *testing.T) {
	assert := assert.New(t)
	db := setupTestDB(t)

	repo := newTestPersonRepository(db)

	found, err := repo.FindByCPF("99999999999")

	assert.NoError(err)
	assert.Nil(found)
}

func TestPersonRepositoryImpl_FindByCPF_WithMultiplePersons(t *testing.T) {
	assert := assert.New(t)
	db := setupTestDB(t)

	repo := newTestPersonRepository(db)

	person1, _ := personModel.NewPerson(
		"Alice Smith",
		"11144477735",
		time.Date(1985, time.March, 15, 0, 0, 0, 0, time.UTC),
		"81987654321",
		"alice@example.com",
	)

	person2, _ := personModel.NewPerson(
		"Bob Johnson",
		"22233344405",
		time.Date(1992, time.July, 22, 0, 0, 0, 0, time.UTC),
		"81998765432",
		"bob@example.com",
	)

	repo.Save(person1)
	repo.Save(person2)

	found, err := repo.FindByCPF("22233344405")

	assert.NoError(err)
	assert.NotNil(found)
	assert.Equal("Bob Johnson", found.Name)
	assert.Equal("22233344405", found.CPF)
}
