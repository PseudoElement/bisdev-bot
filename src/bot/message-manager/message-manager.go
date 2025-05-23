package msg_manager

import (
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	api_module "github.com/pseudoelement/golang-utils/src/api"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type MessageManager struct {
	messages []models.ClientMsg
}

func NewMessageManager() *MessageManager {
	return &MessageManager{
		messages: make([]models.ClientMsg, 0, 200),
	}
}

func (this *MessageManager) SaveNewMsg(update tgbotapi.Update) {
	msg := models.ClientMsg{
		UserName: update.Message.From.UserName,
		Text:     update.Message.Text,
		Id:       time.Now().UnixMicro(),
		IsNew:    true,
	}
	msgToApi := models.JsonClientMsg{
		UserName: update.Message.From.UserName,
		Text:     update.Message.Text,
	}

	go func() {
		_, err := api_module.Post[any](
			"https://rubic-api.exchange/api/v2/buisdev-bot/save-message",
			msgToApi,
			map[string]string{"Content-Type": "application/json"},
		)
		if err != nil {
			log.Println("[SaveNewMsg] error:", err)
		}
	}()

	this.messages = append([]models.ClientMsg{msg}, this.messages...)
}

func (this *MessageManager) GetMessages(count int, newOnly bool) []models.JsonClientMsg {
	params := map[string]string{"count": strconv.Itoa(count)}
	resp, err := api_module.Get[models.GetMessagesResp](
		"https://rubic-api.exchange/api/v2/buisdev-bot/get-messages",
		params,
		map[string]string{"Content-Type": "application/json"},
	)
	if err == nil {
		return resp.Messages
	}

	log.Println("[GetMessages] error:", err)

	messagesLocal := this.getMessagesLocally(count, newOnly)

	return messagesLocal
}

func (this *MessageManager) getMessagesLocally(count int, newOnly bool) []models.JsonClientMsg {
	if count > len(this.messages) {
		count = len(this.messages)
	}

	jsonMessages := make([]models.JsonClientMsg, 0, count)
	checkedMessages := make([]models.ClientMsg, 0, count)

	added := 0
	for i := 0; i < len(this.messages); i++ {
		msg := this.messages[i]

		// stop loop if necessary messages already pushed into slice
		if added >= count {
			break
		}

		if !newOnly || (newOnly && msg.IsNew) {
			jsonMessages = append(jsonMessages, models.JsonClientMsg{
				UserName: msg.UserName,
				Text:     msg.Text,
			})
			checkedMessages = append(checkedMessages, msg)

			msg.IsNew = false
			added++
		}
	}

	go this.deleteMessagesLocallyAfter(checkedMessages, 30)

	return jsonMessages
}

func (this *MessageManager) deleteMessagesLocallyAfter(messages []models.ClientMsg, delayMin int) {
	time.Sleep(time.Duration(delayMin) * time.Minute)

	checkedMsgMap := make(map[int64]models.ClientMsg, len(messages))
	for _, msg := range messages {
		checkedMsgMap[msg.Id] = msg
	}

	notCheckedMessages := make([]models.ClientMsg, 0, len(this.messages)-len(messages))
	for _, msg := range this.messages {
		_, ok := checkedMsgMap[msg.Id]
		if !ok {
			notCheckedMessages = append(notCheckedMessages, msg)
		}
	}

	this.messages = notCheckedMessages
}
