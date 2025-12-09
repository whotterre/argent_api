package dto

type CreateNewUserRequest struct {
	GoogleID string `json:"google_id"`
	Email    string `json:"email"`
	FirstName string    `json:"first_name"`
	LastName string     `json:"last_name"`
}
