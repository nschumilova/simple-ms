package usecase

import (
	"errors"
	"fmt"

	"github.com/nschumilova/simple-ms/pkg/user/domain"
)

type UserRepository interface {
	SelectById(id uint) (*domain.User, error)
	SelectByLogin(login string) (*domain.User, error)
	Insert(params *domain.UserParameters) (*domain.User, error)
	Delete(id uint) (int, error)
	Update(user *domain.User) error
}

type UserUseCase struct {
	userRepository UserRepository
}

func NewUseCase(userRepository UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepository,
	}
}

func (uc *UserUseCase) Create(params *domain.UserParameters) (*domain.User, error) {
	user, error := uc.userRepository.SelectByLogin(params.Login)
	if error != nil && !errors.Is(error, domain.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to create new user: %w", error)
	}
	if user != nil {
		return nil, domain.ErrUserWithSameLoginAlreadyExist
	}
	return uc.userRepository.Insert(params)
}

func (uc *UserUseCase) Delete(userId uint) error {
	rowsAffected, err := uc.userRepository.Delete(userId)
	if err == nil && rowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return err
}

func (uc *UserUseCase) FindById(userId uint) (*domain.User, error) {
	user, error := uc.userRepository.SelectById(userId)
	if error != nil {
		return nil, fmt.Errorf("failed to find user by id %d: %w", userId, error)
	}
	return user, nil
}

func (uc *UserUseCase) Update(userId uint, params *domain.UserParameters) error {
	userById, error := uc.userRepository.SelectById(userId)
	if error != nil {
		return fmt.Errorf("failed to find user by id %d to update: %w", userId, error)
	}
	if userById.Login != params.Login {
		userByLogin, error := uc.userRepository.SelectByLogin(params.Login)
		if error != nil && !errors.Is(error, domain.ErrUserNotFound) {
			return fmt.Errorf("failed to find user by login to update: %w", error)
		}
		if userByLogin!=nil && userByLogin.Id != userId {
			return domain.ErrUserWithSameLoginAlreadyExist
		}
	}
	updatedUser := &domain.User{
		Id:        userId,
		Login:     params.Login,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Phone:     params.Phone,
	}
	return uc.userRepository.Update(updatedUser)
}
