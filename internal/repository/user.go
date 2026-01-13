package repository

import (
	"errors"
	"gosir/internal/database"
	usermodel "gosir/internal/model/user"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

func (r *UserRepository) FindByID(id string) (*usermodel.User, error) {
	var userModel usermodel.User
	err := r.db.Where("id = ?", id).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &UserNotFoundError{ID: id}
		}
		return nil, err
	}
	return &userModel, nil
}

// FindByEmail 通过邮箱查找用户
func (r *UserRepository) FindByEmail(email string) (*usermodel.User, error) {
	var userModel usermodel.User
	err := r.db.Where("email = ?", email).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &UserNotFoundError{ID: email}
		}
		return nil, err
	}
	return &userModel, nil
}

// FindByEmailOrPhone 通过邮箱或手机号查找用户
func (r *UserRepository) FindByEmailOrPhone(account string) (*usermodel.User, error) {
	var userModel usermodel.User
	err := r.db.Where("email = ? OR phone = ?", account, account).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &UserNotFoundError{ID: account}
		}
		return nil, err
	}
	return &userModel, nil
}

func (r *UserRepository) Create(userModel *usermodel.User) (*usermodel.User, error) {
	err := r.db.Create(userModel).Error
	if err != nil {
		return nil, err
	}
	return userModel, nil
}

func (r *UserRepository) FindAll() ([]*usermodel.User, error) {
	var users []*usermodel.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(userModel *usermodel.User) (*usermodel.User, error) {
	err := r.db.Save(userModel).Error
	if err != nil {
		return nil, err
	}
	return userModel, nil
}

func (r *UserRepository) Delete(id string) error {
	err := r.db.Where("id = ?", id).Delete(&usermodel.User{}).Error
	if err != nil {
		return err
	}
	return nil
}

type UserNotFoundError struct {
	ID string
}

func (e *UserNotFoundError) Error() string {
	return "user not found: " + e.ID
}
