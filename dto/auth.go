package dto

type UserSignupDto struct {
	FirstName string `json:"first_name" validate:"required" faker:"first_name"`
	LastName  string `json:"last_name" validate:"required" faker:"last_name"`
	Email     string `json:"email" validate:"required,email" faker:"email"`
	Password  string `json:"password" validate:"required,min=8" faker:"password"`
}

type UserLoginDto struct {
	Email    string `json:"email" validate:"required,email" faker:"email"`
	Password string `json:"password" validate:"required,min=8" faker:"password"`
}

type ResetPasswordDto struct {
	Email string `json:"email" validate:"required,email" faker:"email"`
}
type ConfirmResetPasswordDto struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}
