package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"root/config"
)

var client=&http.Client{Timeout:5*time.Second}

func SendContext(ctx context.Context,message string,cfg *config.Telegram) error {
	payload,err:=json.Marshal(map[string]string{"chat_id":cfg.ChatID,"text":message})
	if err!=nil { return err }
	req,err:=http.NewRequestWithContext(ctx,http.MethodPost,"https://api.telegram.org/bot"+cfg.Token+"/sendMessage",bytes.NewReader(payload))
	if err!=nil { return fmt.Errorf("invalid Telegram request") }
	req.Header.Set("Content-Type","application/json")
	resp,err:=client.Do(req)
	if err!=nil { return fmt.Errorf("Telegram request failed") } // Do not log bot-token URLs.
	defer resp.Body.Close()
	if resp.StatusCode!=http.StatusOK { return fmt.Errorf("Telegram returned HTTP %d",resp.StatusCode) }
	return nil
}
