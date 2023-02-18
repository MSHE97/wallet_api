package models

// Money - денежная суммы в минимальных еденицах (дирамы, копейки, центы и т.д.)
type Money int64

// PaymentCategory - представляет собой категорию, в которой был совершён платёж (cafe, auto, food, drugs, ...)
type PaymentCategory string

// PaymentStatus - представляет собой статус платежа
type PaymentStatus string

type Phone string

// Предопределённые статусы
const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)
