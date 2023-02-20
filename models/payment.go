package models

import (
	"errors"
	"time"
	"wallet/db"
	"wallet/logger"
	"wallet/utils"

	"github.com/jinzhu/gorm"
)

// Payments представляет информацию о платеже
type Payments struct {
	ID          int           `gorm:"column:id; primary_key; auto_increment" json:"id"`
	FromUser    int           `gorm:"column:from_user" json:"from_user"`
	User        Users         `gorm:"-" json:"-"`
	ToAccount   Phone         `gorm:"column:to_account" json:"to_account"`
	ReceiverAcc Accounts      `gorm:"-" json:"-"`
	Amount      Money         `gorm:"column:amount" json:"amount"`
	Type        PaymentType   `gorm:"column:type" json:"type"`
	State       PaymentStatus `gorm:"column:state" json:"state"`
	CreatedAt   time.Time     `gorm:"default: CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt   time.Time     `gorm:"default: null" json:"-"`
}

var (
	ErrPaymentNotFound = errors.New("payment not found ")
	ErrCreatingPayment = errors.New("error creating payment ")
	ErrWrongAmount     = errors.New("wrong amount ")
	ErrWrongReceiveAcc = errors.New("wrong receiver account ")
	ErrOutOfLimit      = errors.New("not enough limit ")
	ErrNotEnBalance    = errors.New("not enough balance ")
	ErrInactiveUser    = errors.New("user inactive ")
)

func (p *Payments) Create(tx *gorm.DB) error {
	p.CreatedAt = time.Now()
	if err := tx.Save(p).Error; err != nil {
		logger.File.Println(ErrCreatingPayment, err)
		return ErrCreatingPayment
	}
	return nil
}

func (p *Payments) Update(tx *gorm.DB) {
	p.UpdatedAt = time.Now()
	if err := tx.Model(p).Update(p).Error; err != nil {
		logger.File.Println("	[WARN] payment update ", p, ". ", err)
	}
}

func (p *Payments) GetByID(ID int64) error {
	if err := db.GetConn().Last(p, ID).Error; err != nil {
		logger.File.Println(ErrPaymentNotFound, "by id =", ID)
		return ErrPaymentNotFound
	}
	return nil
}

func (p *Payments) RunProcessing(tx *gorm.DB) error {
	var (
		xUser, toUser Users
		xAcc, toAcc   Accounts
	)
	xUser.GetByID(p.FromUser)
	toAcc.GetByPhone(p.ToAccount)
	xAcc.GetByID(xUser.AccountId)
	toUser.GetByUUID(toAcc.UserUuid)

	p.State = PaymentStatusInProgress
	// tx.Commit()
	p.Update(tx)
	// TODO: снятие с отправителя, начисление получателю, проверка активности пользователей
	if xAcc.Balance < p.Amount {
		logger.File.Printf("	[CASH-IN] transaction faield. Payment ID %v. %v ", p.ID, ErrNotEnBalance)
		return ErrNotEnBalance
	}
	if !xUser.Active || !toUser.Active {
		logger.File.Printf("	[CASH-IN] transaction faield. Payment ID %v. %v ", p.ID, ErrInactiveUser)
		return ErrInactiveUser
	}
	if toAcc.Balance+p.Amount > Money(utils.Sets.Business.IdentAccLimit) {
		logger.File.Printf("	[CASH-IN] transaction faield. Payment ID %v. %v ", p.ID, ErrOutOfLimit)
		return ErrOutOfLimit
	}

	xAcc.Balance -= p.Amount
	toAcc.Balance += p.Amount
	xAcc.Update(tx)
	toAcc.Update(tx)
	return nil
}
