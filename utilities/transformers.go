package utilities

import (
	"carbon/models"
)

func TransformUsers(users []models.User, fields []string) []models.User {
	transformedUsers := make([]models.User, len(users))

	for index, user := range users {
		switch user.AuthProvider {
		case 1:
			user.AuthenticationProvider = "EMAIL"
		case 2:
			user.AuthenticationProvider = "GOOGLE"
		case 3:
			user.AuthenticationProvider = "GITHUB"
		default:
			user.AuthenticationProvider = "EMAIL"
		}

		transformedUsers[index] = user
	}

	return transformedUsers
}
