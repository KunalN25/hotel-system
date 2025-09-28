package services

import (
	"encoding/json"
	"hotel-system/src/pb"
	"hotel-system/src/store"
	"hotel-system/src/utils"
	"hotel-system/src/validators"
	"net/http"
)

func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest pb.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = validators.ValidateLoginRequest(&loginRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.storageService.GetUserByUsername(loginRequest.Username)
	if err != nil {
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(loginRequest.Password, user.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	response := &pb.LoginOrRegisterResponse{
		Token:   token,
		Message: "Registered successfully",
	}
	sendJsonResponse(w, response)
}

func (s *Service) Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest pb.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = validators.ValidateRegisterRequest(&registerRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	existingUser, err := s.storageService.GetUserByUsername(registerRequest.Username)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	hashedPassword, err := utils.HashPassword(registerRequest.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	newUser := &store.User{
		Username: registerRequest.Username,
		Password: hashedPassword,
	}
	userId, err := s.storageService.AddUser(newUser)
	if err != nil {
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	//	Generate token with the user id
	token, err := utils.GenerateJWT(userId)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	response := &pb.LoginOrRegisterResponse{
		Token:   token,
		Message: "Registered successfully",
	}
	sendJsonResponse(w, response)
}
