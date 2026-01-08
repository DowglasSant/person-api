package ports

import operator "pessoas-api/internal/domain/operator/model"

type OperatorRepository interface {
	Save(operator *operator.Operator) (ID int, err error)
	FindByUsername(username string) (*operator.Operator, error)
	FindByEmail(email string) (*operator.Operator, error)
	FindByID(id int) (*operator.Operator, error)
}
