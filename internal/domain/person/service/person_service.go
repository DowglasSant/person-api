package person

import (
	contract "pessoas-api/internal/contract/person"
	personError "pessoas-api/internal/domain/person/error"
	person "pessoas-api/internal/domain/person/model"
	"pessoas-api/internal/domain/person/ports"
	personUtils "pessoas-api/internal/domain/person/utils"
)

// PersonServiceImpl implements the ports.PersonService interface.
// This is the concrete implementation of the business logic for person operations.
type PersonServiceImpl struct {
	repository ports.PersonRepository
}

// NewPersonService creates a new instance of PersonServiceImpl.
// It returns the implementation as the PersonService interface.
func NewPersonService(repository ports.PersonRepository) ports.PersonService {
	return &PersonServiceImpl{
		repository: repository,
	}
}

func (s *PersonServiceImpl) CreatePerson(newPersonDTO contract.NewPersonDTO) (ID int, err error) {
	person, err := person.NewPerson(
		newPersonDTO.Name,
		newPersonDTO.CPF,
		newPersonDTO.BirthDate,
		newPersonDTO.PhoneNumber,
		newPersonDTO.Email,
	)

	if err != nil {
		return 0, err
	}

	return s.repository.Save(person)
}

func (s *PersonServiceImpl) ListPersons(page, pageSize int, sort, order string) ([]*person.Person, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	if sort == "" {
		sort = "id"
	}
	if order == "" {
		order = "desc"
	}

	return s.repository.FindAll(page, pageSize, sort, order)
}

func (s *PersonServiceImpl) FindPersonByCPF(cpf string) (*person.Person, error) {
	cpfDigits := personUtils.OnlyDigits(cpf)
	return s.repository.FindByCPF(cpfDigits)
}

func (s *PersonServiceImpl) FindPersonByID(id int) (*person.Person, error) {
	return s.repository.FindByID(id)
}

func (s *PersonServiceImpl) UpdatePerson(id int, dto contract.UpdatePersonDTO) error {
	existingPerson, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	if existingPerson == nil {
		return personError.ErrPersonNotFound
	}

	updatedPerson, err := person.NewPerson(
		dto.Name,
		dto.CPF,
		dto.BirthDate,
		dto.PhoneNumber,
		dto.Email,
	)

	if err != nil {
		return err
	}

	updatedPerson.ID = id
	updatedPerson.CreatedAt = existingPerson.CreatedAt

	return s.repository.Update(updatedPerson)
}

func (s *PersonServiceImpl) DeletePerson(id int) error {
	existingPerson, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	if existingPerson == nil {
		return personError.ErrPersonNotFound
	}

	return s.repository.Delete(id)
}
