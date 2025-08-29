package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sancheschris/ecommerce-api/internal/dto"
	"github.com/sancheschris/ecommerce-api/internal/model"
	repo "github.com/sancheschris/ecommerce-api/internal/repository/user"
)

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserDB repo.UserInterface
}

func NewUserHandler(userDB repo.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: userDB,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	u, err := model.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusBadRequest)
		return
	}
	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}