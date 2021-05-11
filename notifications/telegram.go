package notifications

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/rycus86/release-watcher/env"
	"github.com/rycus86/release-watcher/model"
	"github.com/rycus86/release-watcher/transport"
)

// TelegramNotificationManager is a Telegram object that
// the handler receives every time an user interacts with the bot.
type TelegramNotificationManager struct {
	Token   string  `json:"token"`
	ChatID  string  `json:"update_id"`
	Message Message `json:"message"`

	httpClient 	*http.Client
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text     string   `json:"text"`
	Chat     Chat     `json:"chat"`
}

// Chat A Telegram Chat indicates the conversation to which the message belongs.
type Chat struct {
	Id int `json:"id"`
}

func (m *TelegramNotificationManager) initialize() {
	m.Token = env.Lookup("TELEGRAM_BOT_TOKEN", "/var/secrets/telegram", "")
	m.ChatID = env.Lookup("TELEGRAM_CHAT_ID", "/var/secrets/telegram", "")
	//chatID := env.Lookup("TELEGRAM_CHAT_ID", "/var/secrets/telegram", "")
	//m.chatID, err = strconv.Atoi(chatID)
	m.httpClient = &http.Client{
		Timeout:   env.GetTimeout("HTTP_TIMEOUT", "/var/secrets/telegram"),
		Transport: &transport.HttpTransportWithUserAgent{},
	}
}

func (m *TelegramNotificationManager) Close() {
}

// SendNotification sends a text message to the Telegram chat identified by its chat Id
func (m *TelegramNotificationManager) SendNotification(release *model.Release) error {
	var telegramApi = "https://api.telegram.org/bot" + m.Token + "/sendMessage"
	if m.httpClient == nil {
		return errors.New("HTTP client not configured")
	}

	text := fmt.Sprintf("`[New release]` *%s* : %s", release.Project, release.Name)

	if release.URL != "" {
		text = fmt.Sprintf("%s\n%s", text, release.URL)
	}

	response, err := m.httpClient.PostForm(
		telegramApi,
		url.Values{
			"chat_id": 		{m.ChatID},
			"text":    		{text},
			"parse_mode":	{"markdown"},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("error when read http body: %s", err.Error())
		}
	}(response.Body)

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return nil
}
