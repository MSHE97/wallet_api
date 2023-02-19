package models

import "wallet/utils"

type CheckRequest struct {
	ID    int64 `json:"-"`
	Phone Phone `json:"phone"`
}

type TotalsRequest struct {
	ID    int64 `json:"-"`
	Phone Phone `json:"phone"`
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
