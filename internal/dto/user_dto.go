package dto

type CreateUserRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type GetJWTRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type GetJWTResponse struct {
	AcessToken string `json:"access_token"`
}