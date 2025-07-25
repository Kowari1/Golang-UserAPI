package service

import (
	"context"
	"userapi/internal/model"
	"userapi/internal/repository"

	"userapi/internal/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo         repository.UserRepository
	redisService *RedisService
	jwtKey       []byte
}

func NewUserService(repo repository.UserRepository, redisService *RedisService, jwtKey []byte) *UserService {
	return &UserService{
		repo:         repo,
		jwtKey:       jwtKey,
		redisService: redisService,
	}
}

func (s *UserService) Register(user model.User) error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.repo.WithTransaction(func(tx *gorm.DB) error {
		exists, err := s.repo.ExistsByLoginTx(tx, user.Login)
		if err != nil {
			return err
		}

		if exists {
			return &errors.ConflictError{Field: "login", Value: user.Login}
		}

		return tx.Create(&user).Error
	})
}

func (s *UserService) Login(login, password string) (string, error) {
	user, err := s.repo.GetByLogin(login)

	if err != nil {
		return "", err
	}

	if !CheckPassword(user.Password, password) {
		return "", &errors.UnauthorizedError{Reason: "invalid credentials"}
	}

	return GenerateToken(user.ID.String(), user.Admin, user.Login, s.jwtKey)
}

func (s *UserService) GetById(id uuid.UUID) (*model.User, error) {
	user, err := s.repo.GetById(id)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	go func() {
		_ = s.redisService.SetCachedUsers(context.Background(), users)
	}()

	return users, nil
}

func (s *UserService) Delete(id uuid.UUID) error {
	err := s.repo.DeleteWithTransaction(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Update(user model.User) error {

	hashedPassword, err := HashPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return s.repo.UpdateWithTransaction(&user)
}

func (s *UserService) GetByLogin(login string) (*model.User, error) {
	user, err := s.repo.GetByLogin(login)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) EnsureDefaultAdmin() error {
	hasAdmin, err := s.repo.HasAdmin()
	if err != nil {
		return err
	}

	if hasAdmin {
		return nil
	}

	admin := model.User{
		ID:       uuid.New(),
		Login:    "admin",
		Password: "Admin123",
		Name:     "Administrator",
		Gender:   2,
		Admin:    true,
	}

	hashed, err := HashPassword(admin.Password)
	if err != nil {
		return err
	}

	admin.Password = hashed
	admin.CreatedBy = "system"

	return s.repo.Create(&admin)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	return err == nil
}
