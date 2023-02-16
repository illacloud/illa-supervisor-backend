package model

import (
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserStorage struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

func NewUserStorage(db *gorm.DB, logger *zap.SugaredLogger) *UserStorage {
	return &UserStorage{
		logger: logger,
		db:     db,
	}
}

func (d *UserStorage) Create(u *User) (int, error) {
	if err := d.db.Create(u).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}

func (d *UserStorage) RetrieveByID(id int) (*User, error) {
	u := &User{}
	if err := d.db.First(u, id).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *UserStorage) RetrieveByIDs(ids []int) ([]*User, error) {
	users := []*User{}
	if err := d.db.Where("(id) IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (d *UserStorage) RetrieveByUID(uid uuid.UUID) (*User, error) {
	u := &User{}
	if err := d.db.Where("uid = ?", uid).First(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *UserStorage) RetrieveByIDAndUID(id int, uid uuid.UUID) (*User, error) {
	u := &User{}
	if err := d.db.Where("id = ? and uid = ?", id, uid).First(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *UserStorage) RetrieveByEmail(email string) (*User, error) {
	u := &User{}
	if err := d.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *UserStorage) UpdateByID(u *User) error {
	if err := d.db.Model(&User{}).Where("id = ?", u.ID).UpdateColumns(u).Error; err != nil {
		return err
	}
	return nil
}

func (d *UserStorage) DeleteByID(id int) error {
	if err := d.db.Delete(&User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (d *UserStorage) DeleteByUID(uid string) error {
	if err := d.db.Delete(&User{}).Where("uid = ?", uid).Error; err != nil {
		return err
	}
	return nil
}

func (d *UserStorage) ValidateUser(id int, uid uuid.UUID) (bool, error) {
	user, err := d.RetrieveByUID(uid)
	if err != nil {
		return false, err
	}
	if user.ID != id || user.UID != uid {
		return false, errors.New("no such user")
	}
	return true, nil
}
