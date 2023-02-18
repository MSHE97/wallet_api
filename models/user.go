package models

import (
	"errors"
	"time"
	"wallet/db"
	"wallet/logger"

	"github.com/google/uuid"
)

type Users struct {
	ID        int64     `gorm:"column:id; primary_key; auto_increment" json:"id"`
	Uuid      uuid.UUID `gorm:"column:uuid; unique_index; default: null" json:"uuid"`
	Name      string    `gorm:"column:name;" json:"name"`
	Phone     string    `gorm:"column:phone; unique_index" json:"phone"`
	PIN       string    `gorm:"column:pin; default: null" json:"-"`
	AccountId int64     `gorm:"column:account_id" json:"account_id"`
	Active    bool      `gorm:"column:active; default:true" json:"active"`
	BirthDate time.Time `gorm:"column:birth_date; default: null" json:"birth_date"`
	CreatedAt time.Time `gorm:"default: CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time `gorm:"default: null" json:"-"`
}

var (
	ErrUserNotFound = errors.New("user not found ")
	ErrCreatingUser = errors.New("error creating user ")
)

func (u *Users) Create() error {
	u.CreatedAt = time.Now()
	if err := db.GetConn().Save(u).Error; err != nil {
		logger.File.Println(ErrCreatingUser, err)
		return ErrCreatingUser
	}
	return nil
}

func (u *Users) GetByID(ID int64) error {
	if err := db.GetConn().Last(u, ID).Error; err != nil {
		logger.File.Println(ErrUserNotFound, "by id =", ID)
		return ErrUserNotFound
	}
	return nil
}

func (u *Users) GetByPhone(phone string) error {
	if err := db.GetConn().Where(Users{Phone: phone, Active: true}).Last(u).Error; err != nil {
		logger.File.Println(ErrUserNotFound, "phone =", phone)
		return ErrUserNotFound
	}
	return nil
}
