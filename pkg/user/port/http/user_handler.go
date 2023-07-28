package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/nschumilova/simple-ms/pkg/user/domain"
	"go.uber.org/zap"
)

const (
	userIdURLParam = "userID"
)

type UserUsecase interface {
	Create(params *domain.UserParameters) (*domain.User, error)
	FindById(userId uint) (*domain.User, error)
	Update(userId uint, params *domain.UserParameters) error
	Delete(userId uint) error
}

type UserHandler struct {
	uc  UserUsecase
	log *zap.SugaredLogger
}

func NewHandler(useCase UserUsecase, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		uc:  useCase,
		log: logger,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newUser UserBrief
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		h.log.Errorw("failed to parse new user request", "error", err)
		h.error(w, http.StatusBadRequest, "Invalid request")
		return
	}
	createdUser, err := h.uc.Create(toUserParameters(&newUser))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserWithSameLoginAlreadyExist):
			h.error(w, http.StatusConflict, "User already exists")
		default:
			h.log.Errorw("failed to create new user", "newUser", newUser, "error", err)
			h.error(w, http.StatusInternalServerError, "Failed to create user")
		}
		return
	}
	h.success(w, userToView(createdUser))
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, userIdURLParam))
	if err != nil {
		h.error(w, http.StatusBadRequest, "Invalid user id")
		return
	}
	user, err := h.uc.FindById(uint(userId))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			h.error(w, http.StatusNotFound, "User not found")
		default:
			h.log.Errorw("failed to find user by id", "userId", userId, "error", err)
			h.error(w, http.StatusInternalServerError, "Failed to find user by id")
		}
		return
	}
	h.success(w, userToView(user))
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user UserBrief
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.log.Errorw("failed to parse update user request", "error", err)
		h.error(w, http.StatusBadRequest, "Invalid request")
		return
	}
	userId, err := strconv.Atoi(chi.URLParam(r, userIdURLParam))
	if err != nil {
		h.error(w, http.StatusBadRequest, "Invalid user id")
		return
	}
	err = h.uc.Update(uint(userId), toUserParameters(&user))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			h.error(w, http.StatusNotFound, "User not found")
		case errors.Is(err, domain.ErrUserWithSameLoginAlreadyExist):
			h.error(w, http.StatusConflict, "User with new username already exists")
		default:
			h.log.Errorw("failed to update user", "userId", userId, "error", err)
			h.error(w, http.StatusInternalServerError, "Failed to update user")
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, userIdURLParam))
	if err != nil {
		h.error(w, http.StatusBadRequest, "Invalid user id")
		return
	}
	err = h.uc.Delete(uint(userId))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			h.error(w, http.StatusNotFound, "User not found")
		default:
			h.log.Errorw("failed to delete user", "userId", userId, "error", err)
			h.error(w, http.StatusInternalServerError, "Failed to delete user")
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toUserParameters(view *UserBrief) *domain.UserParameters {
	return &domain.UserParameters{
		Login:     view.UserName,
		FirstName: view.FirstName,
		LastName:  view.LastName,
		Email:     view.Email,
		Phone:     view.Phone,
	}
}

func userToView(user *domain.User) *UserFull {
	return &UserFull{
		Id:        user.Id,
		UserName:  user.Login,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}
}

func (h *UserHandler) error(w http.ResponseWriter, code int, message string) {
	view := &ErrorDescription{
		Code:    code,
		Message: message,
	}
	data, err := json.Marshal(view)
	if err != nil {
		h.marshallingError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func (h *UserHandler) success(w http.ResponseWriter, body interface{}) {
	data, err := json.Marshal(body)
	if err != nil {
		h.marshallingError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *UserHandler) marshallingError(w http.ResponseWriter, err error) {
	h.log.Errorw("Can not marshal response to json", "error", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Response building problem."))
}
