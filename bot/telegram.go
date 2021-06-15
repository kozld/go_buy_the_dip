package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type sendMessageReqBody struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func (b *BinanceBot) SendMessage(text string) error {

	reqBody := &sendMessageReqBody{
		ChatID: b.Config.TelegramChannelName,
		Text:   text,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		"https://api.telegram.org/bot"+b.Config.TelegramBotToken+"/"+"sendMessage",
		"application/json",
		bytes.NewBuffer(reqBytes),
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status" + resp.Status)
	}

	return err
}
