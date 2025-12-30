package service

import (
	"myapp/internal/database"
	"myapp/internal/model"
	"myapp/internal/repository"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) CreateUser(name, email, password string) (*model.User, error) {
	user := &model.User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Password:  password,
		Status:    int(model.UserStatusNormal),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.userRepo.Create(user)
}

func (s *UserService) GetAllUsers() ([]*model.User, error) {
	return s.userRepo.FindAll()
}

func (s *UserService) UpdateUser(id, name, email, phone, avatar string, status *int) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	user.Name = name
	user.Email = email
	user.Phone = phone
	user.Avatar = avatar
	if status != nil {
		user.Status = *status
	}
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	return database.DB.AutoMigrate(
		&model.User{},
	)
}
