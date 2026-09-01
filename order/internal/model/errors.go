package model

import "errors"

var (
	// ErrOrderNotFound — заказа с таким UUID нет.
	ErrOrderNotFound = errors.New("order not found")
	// ErrOrderNotPendingPayment — заказ уже не ждёт оплаты.
	ErrOrderNotPendingPayment = errors.New("order is not pending payment")
	// ErrOrderAlreadyPaid — оплаченный заказ нельзя отменить.
	ErrOrderAlreadyPaid = errors.New("order is already paid")
	// ErrPartsNotFound — часть деталей из заказа отсутствует в каталоге.
	ErrPartsNotFound = errors.New("some parts not found")
	// ErrEmptyPartList — заказ без деталей не имеет смысла.
	ErrEmptyPartList = errors.New("part list is empty")
	// ErrUnsupportedPaymentMethod — способ оплаты не поддерживается.
	ErrUnsupportedPaymentMethod = errors.New("unsupported payment method")
	// ErrInventoryUnavailable — каталог деталей недоступен.
	ErrInventoryUnavailable = errors.New("inventory service unavailable")
	// ErrPaymentFailed — платёжный сервис не смог провести оплату.
	ErrPaymentFailed = errors.New("payment failed")
)
