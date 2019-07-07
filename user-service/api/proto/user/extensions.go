package user

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Gen uuid for gorm event lifecycle
func (model *User) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.New()
	return scope.SetColumn("Id", uuid.String())
}
