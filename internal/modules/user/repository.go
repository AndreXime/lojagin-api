package user

import (
	"errors"

	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *User) (*User, error) {
	result := r.db.Create(user) // passa ponteiro
	if result.Error != nil {
		var sqliteErr sqlite3.Error
		if errors.As(result.Error, &sqliteErr) &&
			sqliteErr.Code == sqlite3.ErrConstraint &&
			sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, ErrEmailExists
		}
		return nil, ErrDatabase
	}

	user.Password = "" // nunca expor senha
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrDatabase
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id int64) (*User, error) {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrDatabase
	}
	return &user, nil
}

func (r *UserRepository) GetAllUsers() ([]User, error) {
	var users []User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, ErrDatabase
	}
	return users, nil
}

func (r *UserRepository) IsEmailTaken(email string, excludeID int64) (bool, error) {
	var count int64
	if err := r.db.Model(&User{}).
		Where("email = ? AND id <> ?", email, excludeID).
		Count(&count).Error; err != nil {
		return false, ErrDatabase
	}
	return count > 0, nil
}

func (r *UserRepository) UpdateUser(user *User) error {
	if err := r.db.Save(user).Error; err != nil {
		return ErrDatabase
	}
	return nil
}

func (r *UserRepository) DeleteUser(id int64) error {
	result := r.db.Delete(&User{}, uint(id))
	if result.Error != nil {
		return ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
