package people

import (
	"effectiveMobile/pkg/domain/task"
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

// People represents a person with their details and associated tasks.
// @swagger:model
type People struct {
	ID             int64   `json:"id" form:"id"`
	Name           string  `json:"name" form:"name"`
	Surname        string  `json:"surname" form:"surname"`
	Patronymic     string  `json:"patronymic" form:"patronymic"`
	Address        string  `json:"address" form:"address"`
	Tasks          []int64 `json:"tasksIds" form:"tasksIds"`
	PassportNumber string  `json:"passportNumber" form:"passportNumber"`
	Password       string  `json:"password" form:"password"`
}

// Request represents a request with detailed task information.
// @swagger:model
type Request struct {
	ID             int64       `json:"id" form:"id"`
	Name           string      `json:"name" form:"name"`
	Surname        string      `json:"surname" form:"surname"`
	Patronymic     string      `json:"patronymic" form:"patronymic"`
	Address        string      `json:"address" form:"address"`
	Tasks          []task.Task `json:"tasks" form:"tasks"`
	PassportNumber string      `json:"passportNumber" form:"passportNumber"`
}

// Filter represents a set of criteria for filtering people.
// @swagger:model
type Filter struct {
	ID             *int64   `json:"id" form:"id"`
	Name           *string  `json:"name" form:"name"`
	Surname        *string  `json:"surname" form:"surname"`
	Patronymic     *string  `json:"patronymic" form:"patronymic"`
	Address        *string  `json:"address" form:"address"`
	Tasks          *[]int64 `json:"tasksIds" form:"tasksIds"`
	PassportNumber *string  `json:"passportNumber" form:"passportNumber"`
}

// Info represents information about a person.
// @swagger:model
type Info struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"  validate:"latin-cyrillic" `
	Surname    string `json:"surname"  validate:"latin-cyrillic"`
	Patronymic string `json:"patronymic"  validate:"latin-cyrillic"`
	Address    string `json:"address"  validate:"latin-cyrillic"`
}

// Validate validates the Info struct.
func (info *Info) Validate() error {
	validate := validator.New()

	if err := validate.RegisterValidation("latin-cyrillic", validateLatinCyrillic); err != nil {
		return err
	}
	err := validate.Struct(info)
	if err != nil {
		// Handle validation errors
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		return fmt.Errorf("info validation errors: %s", strings.Join(validationErrors, ", "))
	}

	return nil // No validation errors
}

// Custom validation function for Latin and Cyrillic characters
func validateLatinCyrillic(fl validator.FieldLevel) bool {
	// Regular expression to match Latin and Cyrillic characters
	regex := regexp.MustCompile(`^[\p{Latin}\p{Cyrillic}\s]+$`)
	return regex.MatchString(fl.Field().String())
}

// Registration represents a registration request.
// @swagger:model
type Registration struct {
	PassportNumber string `json:"passportNumber" validate:"required,passportValidate"`
	Password       string `json:"password" validate:"required"`
}

func (r *Registration) Validate() error {
	validate := validator.New()

	if err := validate.RegisterValidation("passportValidate", validatePassportNumber); err != nil {
		return err
	}
	err := validate.Struct(r)
	if err != nil {
		// Handle validation errors
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		return fmt.Errorf("info validation errors: %s", strings.Join(validationErrors, ", "))
	}

	return nil // No validation errors
}

func validatePassportNumber(fl validator.FieldLevel) bool {
	passportRegex := `^\d{4} \d{6}$`
	regex := regexp.MustCompile(passportRegex)
	return regex.MatchString(fl.Field().String())
}

// Pagination represents pagination parameters.
// @swagger:model
type Pagination struct {
	Limit  int `form:"limit" binding:"omitempty,min=1"`
	Offset int `form:"offset" binding:"omitempty,min=1"`
}
