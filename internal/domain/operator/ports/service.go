package ports

type AuthService interface {
	Register(username, email, password string) (operatorID int, err error)
	Login(username, password string) (token string, err error)
}
