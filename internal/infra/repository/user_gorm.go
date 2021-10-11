// Package repository save data
package repository

import (
	"github.com/babydevbr/golang/internal/domain"
	"gorm.io/gorm"
)

// UserGorm repository.
type UserGorm struct {
	db *gorm.DB
}

// NeWUserGorm repository factory.
func NeWUserGorm(db *gorm.DB) *UserGorm {
	return &UserGorm{
		db: db,
	}
}

// Store an user.
func (g *UserGorm) Store(user *domain.User) (*domain.User, error) {
	if dbc := g.db.Create(user); dbc.Error != nil {
		return nil, dbc.Error
	}

	return user, nil
}

// FindOne user.
func (g *UserGorm) FindOne(query *domain.User, args ...string) (*domain.User, error) {
	e := &domain.User{}
	if dbc := g.db.Where(query, args).First(e); dbc.Error != nil {
		return nil, dbc.Error
	}

	return e, nil
}
