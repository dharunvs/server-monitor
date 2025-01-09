package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"root/config"
	"root/logger"
)

type Notification struct {
	Message		string		`json:"message"`
}

var NotificationDataChannel chan Notification

func init () {
	NotificationDataChannel = make(chan Notification)
}

func CreateNotification(message string) Notification{
	return Notification{
		Message: message,
	}
}

func StartNotifier(cfg *config.Telegram,wg *sync.WaitGroup){
	defer wg.Done()
	for notfication := range NotificationDataChannel {
		err := SendNotification(notfication.Message, cfg)
		if err != nil {
			logger.Error(err)
		}
	}
}

func SendNotification(message string, cfg *config.Telegram) error {
	logger.Info("TODO: send notificaiton:", message)

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.Token)

	payload := map[string]string{
		"chat_id": cfg.ChatID,
		"text":    message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	return nil
}