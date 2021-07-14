package config

import (
	"github.com/go-playground/validator"
)

var validatorInstance *validator.Validate

func init() {
	validatorInstance = validator.New()
}

func GetValidator() *validator.Validate {
	return validatorInstance
}
