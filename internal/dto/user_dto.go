package dto

type CreateUserRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID int64 `json:"int"`
	Name string `json:"name"`
	Email string `json:"email"`
}