package models

import (
	"errors"
	"time"
	"wallet/db"
	"wallet/logger"
)

// Payments представляет информацию о платеже
type Payments struct {
	ID        string          `gorm:"column:id; primary_key; auto_increment" json:"id"`
	AccountId int64           `gorm:"column:account_id" json:"account_id"`
	Account   Accounts        `gorm:"-" json:"-"`
	Amount    Money           `gorm:"column:amount" json:"amount"`
	Category  PaymentCategory `gorm:"column:category" json:"category"`
	State     PaymentStatus   `gorm:"column:state" json:"state"`
	CreatedAt time.Time       `gorm:"default: CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time       `gorm:"default: null" json:"-"`
}

var (
	ErrPaymentNotFound = errors.New("payment not found ")
	ErrCreatingPayment = errors.New("error creating payment ")
)

func (u *Payments) Create() error {
	u.CreatedAt = time.Now()
	if err := db.GetConn().Save(u).Error; err != nil {
		logger.File.Println(ErrCreatingPayment, err)
		return ErrCreatingPayment
	}
	return nil
}

func (u *Payments) GetByID(ID int64) error {
	if err := db.GetConn().Last(u, ID).Error; err != nil {
		logger.File.Println(ErrPaymentNotFound, "by id =", ID)
		return ErrPaymentNotFound
	}
	return nil
}
