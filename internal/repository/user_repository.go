package repository

import (
	"errors"
	customErrors "userapi/internal/errors"
	"userapi/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	Create(user *model.User) error
	GetById(id uuid.UUID) (*model.User, error)
	GetAll() ([]model.User, error)
	GetByLogin(login string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uuid.UUID) error
	ExistsByLogin(login string) (bool, error)
	ExistsByLoginTx(tx *gorm.DB, login string) (bool, error)
	HasAdmin() (bool, error)
	UpdateWithTransaction(user *model.User) error
	DeleteWithTransaction(Id uuid.UUID) error
	WithTransaction(fn func(tx *gorm.DB) error) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetById(id uuid.UUID) (*model.User, error) {
	var user model.User

	err := r.db.Where("id = ?", id).First(&user).Error

	if err != nil {
		return wrapNotFound[model.User](err, "User", "id", id.String())
	}

	return &user, nil
}

func (r *userRepository) GetAll() ([]model.User, error) {
	var users []model.User

	err := r.db.Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetByLogin(login string) (*model.User, error) {
	var user model.User

	err := r.db.Where("login = ?", login).First(&user).Error

	if err != nil {
		return wrapNotFound[model.User](err, "User", "login", login)
	}

	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&model.User{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &customErrors.NotFoundError{Entity: "User", Field: "id", Value: id.String()}
	}

	return nil
}

func (r *userRepository) ExistsByLoginTx(tx *gorm.DB, login string) (bool, error) {
	var user model.User

	err := tx.Select("id").Where("login = ?", login).First(&user).Error

	if err == nil {
		return true, nil
	}

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	return false, err
}

func (r *userRepository) ExistsByLogin(login string) (bool, error) {
	var user model.User

	err := r.db.Select("id").Where("login = ?", login).First(&user).Error

	if err == nil {
		return true, nil
	}

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	return false, err
}

func (r *userRepository) UpdateWithTransaction(user *model.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing model.User

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", user.ID).
			First(&existing).Error; err != nil {
			return err
		}

		existing.Name = user.Name
		existing.Login = user.Login
		existing.Password = user.Password
		existing.Gender = user.Gender
		existing.Admin = user.Admin

		if err := tx.Save(&existing).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *userRepository) DeleteWithTransaction(Id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var user model.User

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", Id).
			First(&user).Error; err != nil {
			return wrapNotFoundErr("User", "id", Id.String(), err)
		}

		if err := tx.Delete(&user).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *userRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *userRepository) HasAdmin() (bool, error) {
	var count int64

	err := r.db.Model(&model.User{}).Where("admin = ?", true).Count(&count).Error

	return count > 0, err
}

func wrapNotFound[T any](err error, entity, field, value string) (*T, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &customErrors.NotFoundError{
			Entity: entity,
			Field:  field,
			Value:  value,
		}
	}

	return nil, err
}

func wrapNotFoundErr(entity, field, value string, err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &customErrors.NotFoundError{
			Entity: entity,
			Field:  field,
			Value:  value,
		}
	}

	return err
}
