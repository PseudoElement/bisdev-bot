package query_builder

import (
	"sync"

	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type AdminQueryBuilder struct {
	mu            sync.RWMutex
	queryMessages map[string]models.MessagesReq
}

func NewAdminQueryBuilder() *AdminQueryBuilder {
	return &AdminQueryBuilder{
		queryMessages: make(map[string]models.MessagesReq, 10),
	}
}

func (this *AdminQueryBuilder) SetOldOrNewQueryMsg(adminName string, str string) {
	this.mu.Lock()

	query := models.MessagesReq{}
	if str == consts.SHOW_NEW_MESSAGES {
		query.NewOnly = true
	} else {
		query.NewOnly = false
	}
	this.queryMessages[adminName] = query

	this.mu.Unlock()
}

func (this *AdminQueryBuilder) SetCountOfQueryMsg(adminName string, count int) {
	this.mu.Lock()

	query, ok := this.queryMessages[adminName]
	if ok {
		query.Count = count
		this.queryMessages[adminName] = query
	} else {
		this.queryMessages[adminName] = models.MessagesReq{Count: count, NewOnly: true}
	}

	this.mu.Unlock()
}

func (this *AdminQueryBuilder) GetQueryMsg(adminName string) models.MessagesReq {
	this.mu.RLock()
	query := this.queryMessages[adminName]
	this.mu.RUnlock()

	return query
}
