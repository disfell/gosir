package service

import (
	"gosir/internal/model"
	"gosir/internal/repository"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
	Phone    string `validate:"omitempty,max=20"`
	Avatar   string `validate:"omitempty,max=500"`
	Status   *int   `validate:"omitempty,oneof=1 2"`
}

func (s *UserService) CreateUser(req *CreateUserRequest) (*model.User, error) {
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 设置默认状态
	status := int(model.UserStatusNormal)
	if req.Status != nil {
		status = *req.Status
	}

	user := &model.User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Phone:     req.Phone,
		Avatar:    req.Avatar,
		Status:    status,
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
