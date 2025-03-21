package dto

import "github.com/google/uuid"

type VerifyUserDto struct {
	Uuid uuid.UUID `param:"uuid"`
}

type UpdateUserDto struct {
	Uuid      uuid.UUID `param:"uuid"`
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
}
