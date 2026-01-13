package user

import (
	usermodel "gosir/internal/model/user"
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

func (s *UserService) GetUserByID(id string) (*usermodel.User, error) {
	return s.userRepo.FindByID(id)
}

// GetUserByEmail 通过邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*usermodel.User, error) {
	return s.userRepo.FindByEmail(email)
}

// GetUserByEmailOrPhone 通过邮箱或手机号获取用户
func (s *UserService) GetUserByEmailOrPhone(account string) (*usermodel.User, error) {
	return s.userRepo.FindByEmailOrPhone(account)
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

func (s *UserService) CreateUser(req *CreateUserRequest) (*usermodel.User, error) {
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 设置默认状态
	status := int(usermodel.UserStatusNormal)
	if req.Status != nil {
		status = *req.Status
	}

	newUser := &usermodel.User{
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
	return s.userRepo.Create(newUser)
}

func (s *UserService) GetAllUsers() ([]*usermodel.User, error) {
	return s.userRepo.FindAll()
}

func (s *UserService) UpdateUser(id, name, email, phone, avatar string, status *int) (*usermodel.User, error) {
	userModel, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	userModel.Name = name
	userModel.Email = email
	userModel.Phone = phone
	userModel.Avatar = avatar
	if status != nil {
		userModel.Status = *status
	}
	userModel.UpdatedAt = time.Now()
	return s.userRepo.Update(userModel)
}

func (s *UserService) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}
