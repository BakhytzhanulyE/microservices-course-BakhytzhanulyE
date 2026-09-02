// Package notify собирает текст уведомления и отправляет его пользователю.
package notify

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/model"
	def "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/service"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

var _ def.NotifyService = (*service)(nil)

// templatesFS — шаблоны сообщений лежат рядом с кодом и вшиты в бинарник.
//
//go:embed templates/*.tmpl
var templatesFS embed.FS

// templates разбираем один раз на старте: парсить их на каждое событие незачем.
var templates = template.Must(template.ParseFS(templatesFS, "templates/*.tmpl"))

// Sender умеет доставить готовый текст пользователю.
type Sender interface {
	SendMessage(ctx context.Context, text string) error
}

type service struct {
	sender Sender
}

// NewService создаёт сервис уведомлений поверх отправщика сообщений.
func NewService(sender Sender) *service {
	return &service{sender: sender}
}

// OrderPaid отправляет уведомление об оплате заказа.
func (s *service) OrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	text, err := render("order_paid.tmpl", event)
	if err != nil {
		logger.Error(ctx, "Не удалось собрать текст уведомления", zap.Error(err))
		return err
	}

	return s.sender.SendMessage(ctx, text)
}

// ShipAssembled отправляет уведомление о собранном корабле.
func (s *service) ShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	text, err := render("ship_assembled.tmpl", event)
	if err != nil {
		logger.Error(ctx, "Не удалось собрать текст уведомления", zap.Error(err))
		return err
	}

	return s.sender.SendMessage(ctx, text)
}

func render(name string, data any) (string, error) {
	var buf bytes.Buffer

	if err := templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
