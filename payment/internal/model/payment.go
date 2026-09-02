// Package model описывает доменные сущности оплаты.
package model

import "errors"

// PaymentMethod — способ оплаты.
type PaymentMethod int

const (
	// PaymentMethodUnspecified — способ оплаты не указан.
	PaymentMethodUnspecified PaymentMethod = iota
	// PaymentMethodCard — банковская карта.
	PaymentMethodCard
	// PaymentMethodSBP — система быстрых платежей.
	PaymentMethodSBP
	// PaymentMethodCreditCard — кредитная карта.
	PaymentMethodCreditCard
	// PaymentMethodInvestorMoney — деньги инвестора.
	PaymentMethodInvestorMoney
)

// ErrUnsupportedPaymentMethod — способ оплаты не поддерживается.
var ErrUnsupportedPaymentMethod = errors.New("unsupported payment method")

// Valid говорит, поддерживается ли способ оплаты.
func (m PaymentMethod) Valid() bool {
	switch m {
	case PaymentMethodCard, PaymentMethodSBP, PaymentMethodCreditCard, PaymentMethodInvestorMoney:
		return true
	case PaymentMethodUnspecified:
		return false
	default:
		return false
	}
}

// String — читаемое имя способа оплаты, оно уезжает в событие OrderPaid.
func (m PaymentMethod) String() string {
	switch m {
	case PaymentMethodCard:
		return "CARD"
	case PaymentMethodSBP:
		return "SBP"
	case PaymentMethodCreditCard:
		return "CREDIT_CARD"
	case PaymentMethodInvestorMoney:
		return "INVESTOR_MONEY"
	case PaymentMethodUnspecified:
		return "UNSPECIFIED"
	default:
		return "UNSPECIFIED"
	}
}

// Payment — запрос на оплату заказа.
type Payment struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod PaymentMethod
}
