package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/abrarr21/test/internal/middlewares"
	"github.com/abrarr21/test/internal/models"
	"github.com/abrarr21/test/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("failed to decode the input: ", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		log.Printf("validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user.ID = bson.NewObjectID()
	user.Email = strings.ToLower(user.Email)
	user.Password = string(hashed)

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

	token, err := utils.GenerateToken(user.ID.Hex(), user.Email, h.Cfg.JWT.Secret)
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		utils.JSON(w, http.StatusInternalServerError, "failed to generate token", nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "accessToken",
		Value:  token,
		MaxAge: 15 * 60,
	})

	response := models.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Address: user.Address,
	}

	if err := utils.JSON(w, http.StatusCreated, "User created successfully", response); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSON(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.JSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	input.Email = strings.ToLower(input.Email)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	err := h.DB.Users.FindOne(ctx, bson.D{{Key: "email", Value: input.Email}}).Decode(&user)
	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, "invalid email or password", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.JSON(w, http.StatusUnauthorized, "invalid email or password", nil)
		return
	}

	tokenString, err := utils.GenerateToken(user.ID.Hex(), user.Email, h.Cfg.JWT.Secret)
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		utils.JSON(w, http.StatusInternalServerError, "failed to generate token", nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "accessToken",
		Value:  tokenString,
		MaxAge: 15 * 60,
	})

	utils.JSON(w, http.StatusOK, "login successful", map[string]string{
		"token": tokenString,
	})
}

func (h *Handlers) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Write([]byte(userID))
}
