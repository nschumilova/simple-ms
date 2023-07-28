package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/nschumilova/simple-ms/pkg/user/domain"
	"gorm.io/gorm"
)

type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement:true"`
	Login        string `gorm:"unique:true"`
	FirstName    string
	LastName     string
	Email        sql.NullString
	Phone        sql.NullString
	CreatedAtUtc time.Time
	UpdatedAtUtc time.Time
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(database *gorm.DB) *UserRepository {
	return &UserRepository{
		db: database,
	}
}

func (r *UserRepository) Delete(id uint) (int,error) {
	result := r.db.Delete(&User{}, id)
	
	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete user by id %d: %w", id, result.Error)
	}
	return int(result.RowsAffected), nil
}

func (r *UserRepository) Insert(params *domain.UserParameters) (*domain.User, error) {
	user := createParamsToEntity(params)
	if result := r.db.Create(user); result.Error != nil {
		return nil, fmt.Errorf("failed to insert user: %w", result.Error)
	}
	return entityToDomain(user), nil
}

func (r *UserRepository) SelectById(id uint) (*domain.User, error) {
	var user User
	result := r.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user by id %d: %w", id, result.Error)
	}
	return entityToDomain(&user), nil
}

func (r *UserRepository) SelectByLogin(login string) (*domain.User, error) {
	var user User
	result := r.db.Where(&User{Login: login}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user by login %v: %w", login, result.Error)
	}
	return entityToDomain(&user), nil
}

func (r *UserRepository) Update(user *domain.User) error {
	entity := domainToEntity(user)
	result := r.db.Save(entity)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.ErrUserNotFound
	}
	if result.Error != nil {
		return fmt.Errorf("failed to update user by id %d: %w", user.Id, result.Error)
	}
	return nil
}

func entityToDomain(entity *User) *domain.User {
	return &domain.User{
		Id:        entity.ID,
		Login:     entity.Login,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Email:     fromNullString(entity.Email),
		Phone:     fromNullString(entity.Phone),
	}
}
func domainToEntity(domain *domain.User) *User {
	return &User{
		ID:           domain.Id,
		Login:        domain.Login,
		FirstName:    domain.FirstName,
		LastName:     domain.LastName,
		Email:        toNullString(domain.Email),
		Phone:        toNullString(domain.Phone),
		UpdatedAtUtc: time.Now().UTC(),
	}
}

func createParamsToEntity(params *domain.UserParameters) *User {
	now := time.Now().UTC()
	return &User{
		Login:        params.Login,
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Email:        toNullString(params.Email),
		Phone:        toNullString(params.Phone),
		CreatedAtUtc: now,
		UpdatedAtUtc: now,
	}
}

func toNullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func fromNullString(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}
