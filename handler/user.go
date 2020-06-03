package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/mochisuna/load-test-sample/domain"
)

type displayRequest struct {
	ID        domain.UserID `json:"id"`
	Name      string        `json:"name"`
	SecretKey string        `json:"secret_key"`
	RequestID string        `json:"request_id"`
}

type userResponse struct {
	ID        domain.UserID `json:"id"`
	Name      string        `json:"name"`
	SecretKey string        `json:"secret_key"`
	RequestID string        `json:"request_id"`
	CreatedAt int           `json:"created_at"`
	UpdatedAt int           `json:"updated_at"`
}

func (s *Server) referUser(w http.ResponseWriter, r *http.Request) {
	log.Println("referUser")
	ctx := r.Context()
	requestID := middleware.GetReqID(ctx)
	userID := chi.URLParam(r, "userID")
	user, err := s.UserService.Refer(domain.UserID(userID))
	if err != nil {
		log.Println("error reason: " + err.Error())
		switch {
		case err == sql.ErrNoRows:
			rendering.JSON(w, http.StatusNotFound, &responseError{
				RequestID: requestID,
				Reason:    "User: " + userID + " is not found.",
				Error:     nil,
			})
		default:
			rendering.JSON(w, http.StatusInternalServerError, &responseError{
				RequestID: requestID,
				Reason:    "Failed UserService.Refer.",
				Error:     err.Error(),
			})
		}
		return
	}
	rendering.JSON(w, http.StatusOK, &userResponse{
		RequestID: requestID,
		ID:        user.ID,
		Name:      user.Name,
		SecretKey: user.SecretKey,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (s *Server) registerUser(w http.ResponseWriter, r *http.Request) {
	log.Println("registerUser")
	ctx := r.Context()
	requestID := middleware.GetReqID(ctx)
	user := &domain.User{}
	if err := s.UserService.Register(ctx, user); err != nil {
		log.Println("error reason: " + err.Error())
		rendering.JSON(w, http.StatusInternalServerError, &responseError{
			RequestID: requestID,
			Reason:    "Failed UserService.Register.",
			Error:     err.Error(),
		})
		return
	}
	rendering.JSON(w, http.StatusCreated, &userResponse{
		RequestID: requestID,
		ID:        user.ID,
		Name:      user.Name,
		SecretKey: user.SecretKey,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (s *Server) displayUser(w http.ResponseWriter, r *http.Request) {
	log.Println("displayUser")
	ctx := r.Context()
	requestID := middleware.GetReqID(ctx)

	// decode request body
	var dr displayRequest
	if err := json.NewDecoder(r.Body).Decode(&dr); err != nil {
		log.Println("error reason: " + err.Error())
		rendering.JSON(w, http.StatusBadRequest, &responseError{
			RequestID: requestID,
			Reason:    "Failed request body decode.",
			Error:     err.Error(),
		})
		return
	}

	// validation
	if err := validate.Struct(dr); err != nil {
		log.Println("error reason: " + err.Error())
		rendering.JSON(w, http.StatusBadRequest, &responseError{
			RequestID: requestID,
			Reason:    "Failed request body validate.",
			Error:     err.Error(),
		})
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("location", fmt.Sprintf("%s?user_id=%s&user_name=%s&secret_key=%s", s.RedirectURL, dr.ID, dr.Name, dr.SecretKey))
	w.WriteHeader(http.StatusMovedPermanently)
}
