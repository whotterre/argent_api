package dto

type CreateNewUserRequest struct {
	GoogleID string `json:"google_id"`
	Email string `json:"email"`
	FullName string `json:"full_name"`
}