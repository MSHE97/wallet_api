package models

import (
	"time"
	"wallet/db"
	"wallet/logger"
	"wallet/utils"

	"github.com/jinzhu/gorm"
)

type CheckRequest struct {
	ID    int64 `json:"-"`
	Phone Phone `json:"phone"`
}

type TotalsRequest struct {
	ID   int64  `json:"-"`
	From string `json:"from"`
	To   string `json:"to"`
}

type BalanceRequest struct {
	ID    int64 `json:"-"`
	Phone Phone `json:"phone"`
}

type PaymentRequest struct {
	ID        int64       `json:"-"`
	Amount    Money       `json:"amount"`
	ToAccount Phone       `json:"phone"`
	Type      PaymentType `json:"type"`
}

func (pay *PaymentRequest) PreCheckSender(xUserId int) error {
	var user Users
	if xUserId == 0 {
		return ErrUserNotFound
	}
	if err := user.GetByID(xUserId); err != nil {
		return err
	}
	if !user.Active {
		return ErrInactiveUser
	}
	return nil
}

func (pay *PaymentRequest) PreCheckReceiver() error {
	var account Accounts
	bKeys := utils.Sets.Business
	if err := account.GetByPhone(pay.ToAccount); err != nil || account.ID == 0 {
		return err
	}

	if err := CheckUser(account); err != nil {
		return err
	}

	switch account.Identified {
	case true:
		if pay.Amount+account.Balance > Money(bKeys.IdentAccLimit) {
			return ErrOutOfLimit
		}
	default:
		if pay.Amount+account.Balance > Money(bKeys.SimpAccLimit) {
			return ErrOutOfLimit
		}
	}
	return nil
}

type PayTotals struct {
	Name          string  `json:"name" gorm:"-"`
	Count         int64   `json:"total_count" gorm:"column:total_count"`
	Amount        float64 `json:"total_amount" gorm:"column:total_amount"`
	Refunds       int64   `json:"refunded" gorm:"column:refunds"`
	RefundsAmount float64 `json:"refunds_amount" gorm:"column:refunds_amount"`
}

func (r *TotalsRequest) GetTotals() (totals PayTotals, err error) {
	totals.Name = "cash-in"
	from, err := time.Parse(time.RFC3339, r.From)
	if err != nil {
		logger.File.Println(" [TOTAL] data parse error ", err)
		return
	}
	to, err := time.Parse(time.RFC3339, r.To)
	if err != nil {
		logger.File.Println(" [TOTAL] data parse error ", err)
		return
	}
	// time validation
	if from.Unix() <= 0 || to.Unix() <= 0 {
		from = TruncateToMounth(time.Now())
		to = TruncateToDay(time.Now())
	}

	if err := db.GetConn().Raw(`
	SELECT t1.total_count, total_amount, t2.refunds, refunds_amount 
	FROM (
			SELECT COUNT(*) as total_count, SUM(payments.amount), id as total_amount 
			FROM public.payments WHERE from_user is not null and to_account is not null and state = ? and created_at between coalesce(?,created_at) and coalesce(?,created_at)
		 ) as t1
	LEFT JOIN (
			SELECT COUNT(*) as refunds, SUM(payments.amount), id as refunds_amount
			FROM public.transaction WHERE from_user is not null and to_account is not null and state = ? and created_at between coalesce(?,created_at) and coalesce(?,created_at)
		) as t2
		on t1.id = t2.id`,
		PaymentStatusOk, from, to, PaymentStatusRefund, from, to).Scan(&totals).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return totals, nil
		} else {
			logger.File.Println("Totals Query error: ", err)
			return PayTotals{}, err
		}
	}
	return
}

func TruncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
}

func TruncateToMounth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, t.Location())
}
