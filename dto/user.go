package dto

import "github.com/google/uuid"

type VerifyUserDto struct {
	Uuid uuid.UUID `param:"uuid" faker:"uuid_hyphenated"`
}

type UpdateUserDto struct {
	Uuid               uuid.UUID `param:"uuid"`
	FirstName          string    `form:"first_name"`
	LastName           string    `form:"last_name"`
	CurrentProjectUuid uuid.UUID `form:"current_project_uuid"`
}
