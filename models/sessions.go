package models

import (
	"errors"
	"time"
	"wallet/db"
	"wallet/logger"
)

type Sessions struct {
	ID        int64      `gorm:"column:id; primary_key; auto_increment" json:"id"`
	UserID    int        `gorm:"column:user_id" json:"user_id"`
	IPAddress string     `gorm:"column:ip_address; default:null" json:"ip_address"`
	LoginAt   time.Time  `gorm:"default: CURRENT_TIMESTAMP" json:"-"`
	LogoutAt  *time.Time `gorm:"default: null" json:"-"`
	CreatedAt time.Time  `gorm:"default: CURRENT_TIMESTAMP" json:"-"`
}

var (
	ErrLogout          = errors.New("user session logout error ")
	ErrSesionNotFound  = errors.New("user session not found ")
	ErrCreatingSession = errors.New("user session create error ")
)

func (s *Sessions) Create() error {
	s.CreatedAt = time.Now()
	if err := db.GetConn().Save(s).Error; err != nil {
		logger.File.Printf("%v id = %v user_id = %v. %v", ErrCreatingSession, s.ID, s.UserID, err)
		return ErrCreatingSession
	}
	return nil
}

func (s *Sessions) Find() error {
	var user Users
	if err := user.GetByID(s.UserID); err != nil {
		return err
	}
	if err := db.GetConn().Where("logout_at is null and user_id = ? and account_id = ?", s.UserID, user.AccountId).Last(s, s.ID).Error; err != nil {
		logger.File.Printf("%v id = %v user_id = %v. %v", ErrSesionNotFound, s.ID, s.UserID, err)
		return ErrSesionNotFound
	}
	return nil
}

func (s *Sessions) Logout() error {
	now := time.Now()
	if err := db.GetConn().Model(s).Update(Sessions{LogoutAt: &now}).Error; err != nil {
		logger.File.Printf("%v id = %v user_id = %v. %v", ErrLogout, s.ID, s.UserID, err)
		return ErrLogout
	}
	return nil
}
