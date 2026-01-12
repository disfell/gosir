package repository

import (
	"errors"
	"gosir/internal/database"
	"gosir/internal/model"

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

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &UserNotFoundError{ID: id}
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail 通过邮箱查找用户
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &UserNotFoundError{ID: email}
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmailOrPhone 通过邮箱或手机号查找用户
func (r *UserRepository) FindByEmailOrPhone(account string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? OR phone = ?", account, account).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &UserNotFoundError{ID: account}
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *model.User) (*model.User, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindAll() ([]*model.User, error) {
	var users []*model.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(user *model.User) (*model.User, error) {
	err := r.db.Save(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Delete(id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.User{}).Error
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
