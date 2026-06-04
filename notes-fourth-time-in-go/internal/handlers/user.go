package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/abrarr21/notes-in-golang/internal/middlewares"
	"github.com/abrarr21/notes-in-golang/internal/models"
	"github.com/abrarr21/notes-in-golang/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var input models.RegisterUserRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("failed to decode request body: %v", err)
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid input", nil)
		return
	}

	if err := utils.Validate(input); err != nil {
		log.Printf("failed to validate: %v", err)
		utils.ResponseJSON(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash the password: %v", err)
		http.Error(w, "failed to hash the Passwrod", http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:       bson.NewObjectID(),
		Name:     strings.TrimSpace(input.Name),
		Email:    strings.ToLower(input.Email),
		Password: string(hashed),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := h.DB.Users.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			http.Error(w, "email already exists", http.StatusConflict)
			return
		}
		log.Printf("failed to insert user: %v", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	log.Println("User created: ", result.InsertedID)

	token, err := utils.GenerateToken(user.ID.Hex(), user.Email, h.Cfg.JWT.JWT_SECRET, h.Cfg.JWT.AccessTokenTTL)
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		utils.ResponseJSON(w, http.StatusInternalServerError, "failed to generate token", nil)
		return
	}

	response := models.UserResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   15 * 60,
	})

	if err := utils.ResponseJSON(w, http.StatusCreated, "user created successfully", response); err != nil {
		log.Println("failed to encode respones")
	}
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var input models.LoginUserRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("failed to decode input: %v", err)
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid input", nil)
		return
	}

	if err := utils.Validate(input); err != nil {
		log.Printf("failed to validate: %v", err)
		utils.ResponseJSON(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	input.Email = strings.ToLower(input.Email)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var user models.User

	err := h.DB.Users.FindOne(ctx, bson.D{{Key: "email", Value: input.Email}}).Decode(&user)
	if err != nil {
		utils.ResponseJSON(w, http.StatusUnauthorized, "invalid email or password", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.ResponseJSON(w, http.StatusUnauthorized, "invalid email or password", nil)
		return
	}

	token, err := utils.GenerateToken(user.ID.Hex(), user.Email, h.Cfg.JWT.JWT_SECRET, h.Cfg.JWT.AccessTokenTTL)
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		utils.ResponseJSON(w, http.StatusInternalServerError, "failed to generate token", nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    token,
		Path:     "/",
		MaxAge:   15 * 60,
		HttpOnly: true,
	})

	response := models.UserResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
	}

	if err := utils.ResponseJSON(w, http.StatusOK, "log in successfully", response); err != nil {
		log.Println("failed to encode response")
	}
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	emailId, ok := middlewares.GetEmailId(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := h.DB.Users.FindOne(ctx, bson.D{{Key: "email", Value: emailId}}).Decode(&user)
	if err != nil {
		utils.ResponseJSON(w, http.StatusOK, "couldn't find user", nil)
		return
	}

	response := models.UserResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
	}

	if err := utils.ResponseJSON(w, http.StatusOK, "user fetched", response); err != nil {
		log.Println("error encoding response")
	}
}

func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	if err := utils.ResponseJSON(w, http.StatusOK, "log out successfully", nil); err != nil {
		log.Println("error encoding resonse")
	}
}
