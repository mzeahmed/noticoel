package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mzeahmed/noticoel/internal/notifier"
)

type Notifier struct {
	client *http.Client
	config Config
}

type sendMessageRequest struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func New(cfg Config) *Notifier {
	return &Notifier{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (n *Notifier) Name() string {
	return "telegram"
}

func (n *Notifier) Notify(ctx context.Context, msg notifier.Message) notifier.Result {

	body := sendMessageRequest{
		ChatID: n.config.ChatID,
		Text: fmt.Sprintf(
			"[%s] %s\n\n%s",
			msg.Severity,
			msg.Title,
			msg.Message,
		),
	}

	payload, _ := json.Marshal(body)

	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		n.config.BotToken,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(payload),
	)

	if err != nil {
		return notifier.Result{
			Notifier: n.Name(),
			Success:  false,
			Error:    err,
			SentAt:   time.Now(),
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)

	if err != nil {
		return notifier.Result{
			Notifier: n.Name(),
			Success:  false,
			Error:    err,
			SentAt:   time.Now(),
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return notifier.Result{
			Notifier: n.Name(),
			Success:  false,
			Message:  resp.Status,
			SentAt:   time.Now(),
		}
	}

	return notifier.Result{
		Notifier: n.Name(),
		Success:  true,
		SentAt:   time.Now(),
	}
}
