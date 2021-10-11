// Package gateway ...
package gateway

import (
	"github.com/babydevbr/golang/internal/domain"
	"github.com/google/uuid"
)

// Auth ....
type Auth struct{}

// GenerateToken ....
func (a *Auth) GenerateToken(user *domain.User) (string, error) {
	return uuid.NewString(), nil
}
