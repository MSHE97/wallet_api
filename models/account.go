package models

import (
	"errors"
	"time"
	"wallet/db"
	"wallet/logger"

	"github.com/google/uuid"
)

type Accounts struct {
	ID           int64     `gorm:"column:id; primary_key; auto_increment" json:"id"`
	Phone        Phone     `gorm:"column:phone" json:"phone"`
	UserUuid     uuid.UUID `gorm:"column:user_uuid" json:"user_uuid"`
	User         Users     `gorm:"-" json:"-"`
	Identified   bool      `gorm:"column:identified" json:"identified"`
	IdentifiedAt time.Time `gorm:"column:default: null" json:"-"`
	Balance      Money     `gorm:"column:balance" json:"balance"`
	CreatedAt    time.Time `gorm:"default: CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt    time.Time `gorm:"default: null" json:"-"`
}

var (
	ErrCreatingAcc = errors.New("creating account error ")
	ErrAccNotFound = errors.New("account not found ")
	ErrShortName   = errors.New("too short name ")
	ErrEmptyUuid   = errors.New("empty uuid ")
)

func (a *Accounts) Create() error {
	a.CreatedAt = time.Now()
	if err := db.GetConn().Save(a).Error; err != nil {
		logger.File.Println(ErrCreatingAcc, err)
		return ErrCreatingAcc
	}
	return nil
}

func (a *Accounts) GetByID(ID int64) error {
	if err := db.GetConn().Last(a, ID).Error; err != nil {
		logger.File.Println(ErrAccNotFound, "id = ", ID)
		return ErrAccNotFound
	}
	return nil
}

func (a *Accounts) GetByPhone(phone string) error {
	if err := db.GetConn().Where(Accounts{Phone: Phone(phone)}).Last(a).Error; err != nil {
		logger.File.Println(ErrAccNotFound, "phone =", phone)
		return ErrAccNotFound
	}
	return nil
}

func (a *Accounts) Identify(user Users) error {
	if len(user.Name) < 1 {
		logger.File.Println(ErrShortName, "user_id = ", user.ID)
		return ErrShortName
	}
	if len(user.Uuid) < 1 {
		logger.File.Println(ErrEmptyUuid, "user_id = ", user.ID)
		return ErrEmptyUuid
	}
	a.UserUuid = user.Uuid
	a.Identified = true
	a.IdentifiedAt = time.Now()
	a.UpdatedAt = time.Now()
	logger.File.Printf("	[IDENTIFY] user id = %v", user.ID)
	return nil
}
