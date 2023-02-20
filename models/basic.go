package models

// Money - денежная суммы в минимальных еденицах (дирамы, копейки, центы и т.д.)
type Money int

// PaymentCategory - представляет собой категорию, в которой был совершён платёж (cafe, auto, food, drugs, ...)
type PaymentCategory string

// PaymentStatus - представляет собой статус платежа
type PaymentStatus string

type PaymentType string

type Phone string

// Предопределённые статусы
const (
	PaymentStatusOk         PaymentStatus = "COMPLETE"
	PaymentStatusRefund     PaymentStatus = "REFUND"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
	PaymentStatusInSaved    PaymentStatus = "SAVED"
)

const (
	CasIn    PaymentType = "cash_in"
	Transfer PaymentType = "smart_transfer"
	Qr       PaymentType = "qr"
)
