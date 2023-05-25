package database

import (
	"errors"

	"gorm.io/gorm"
)

// https://github.com/go-gorm/gorm/blob/master/errors.go
func CheckDatabaseError(err error) string {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "Element does not exist."
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return "Element already exists."
	}

	return err.Error()
}
