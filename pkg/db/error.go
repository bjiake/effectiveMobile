package db

import "errors"

// Обозначение ошибок
var (
	ErrMigrate           = errors.New("migration failed")
	ErrDuplicate         = errors.New("record already exists")
	ErrNotExist          = errors.New("row does not exist")
	ErrUpdateFailed      = errors.New("update failed")
	ErrDeleteFailed      = errors.New("delete failed")
	ErrTasks             = errors.New("get tasks failed")
	ErrParamNotFound     = errors.New("param not found")
	ErrValidate          = errors.New("validate failed")
	ErrPassportSerie     = errors.New("passport serie not valid")
	ErrPassportNumber    = errors.New("passport number not valid")
	ErrTimeInvalidFormat = errors.New("invalid time format")
)
