// Package app use cases
package app

import "github.com/babydevbr/golang/internal/domain"

type storager interface {
	Store(user *domain.User) (*domain.User, error)
	FindOne(query *domain.User, args ...string) (*domain.User, error)
}

type encrypter interface {
	Encrypt(s string) string
	Compare(hash, s string) error
}

type checker interface {
	Struct(s interface{}) error
}

// UserUseCase user auth uses case.
type UserUseCase struct {
	repository storager
	encrypter  encrypter
	validate   checker
}

// NewUserUseCase factory.
func NewUserUseCase(s storager, e encrypter, v checker) *UserUseCase {
	return &UserUseCase{
		repository: s,
		encrypter:  e,
		validate:   v,
	}
}

// SignUp create a new user.
func (u *UserUseCase) SignUp(user *domain.User) (*domain.User, error) {
	if err := u.validate.Struct(user); err != nil {
		return nil, ErrInvalid
	}

	user.Password = u.encrypter.Encrypt(user.Password)

	newUser, err := u.repository.Store(user)
	if err != nil {
		return nil, ErrOnSave
	}

	newUser.Password = ""

	return newUser, nil
}

// GetByEmail Get User By Email.
func (u *UserUseCase) GetByEmail(email string) (*domain.User, error) {
	e := &domain.User{
		Email: email,
	}

	user, err := u.repository.FindOne(e, "email")
	if err != nil {
		return nil, ErrInvalid
	}

	return user, nil
}
