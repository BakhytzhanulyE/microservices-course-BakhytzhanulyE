// Package telegram — клиент Telegram Bot API.
package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

const (
	apiURL         = "https://api.telegram.org"
	requestTimeout = 10 * time.Second
)

// Client отправляет сообщения в Telegram.
type Client struct {
	httpClient *http.Client
	token      string
	chatID     string
}

// NewClient создаёт клиента. Пустой token означает «бот не настроен»:
// сервис в этом случае просто пишет уведомления в лог.
func NewClient(token, chatID string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: requestTimeout},
		token:      token,
		chatID:     chatID,
	}
}

type sendMessageRequest struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

// SendMessage отправляет текст в настроенный чат.
func (c *Client) SendMessage(ctx context.Context, text string) error {
	if c.token == "" || c.chatID == "" {
		logger.Info(ctx, "📭 Telegram не настроен, уведомление только в логе", zap.String("text", text))
		return nil
	}

	payload, err := json.Marshal(sendMessageRequest{ChatID: c.chatID, Text: text})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", apiURL, c.token)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.Error(ctx, "Не удалось закрыть тело ответа Telegram", zap.Error(closeErr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram вернул статус %d", resp.StatusCode)
	}

	return nil
}
