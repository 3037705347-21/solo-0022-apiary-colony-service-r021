package app

import (
	"net/http"
	"time"

	"example.com/apiary-colony-service/internal/apiary"
	"example.com/apiary-colony-service/internal/httpapi"
	"example.com/apiary-colony-service/internal/inspection"
	"example.com/apiary-colony-service/internal/queen"
	"example.com/apiary-colony-service/internal/status"
	"example.com/apiary-colony-service/internal/store"
)

// Application 把内存仓储、领域服务与 HTTP 入口组装在一起。
type Application struct {
	Config Config
	Store  *store.Memory
	Server *httpapi.Server
}

func New(config Config) *Application {
	return NewWithClock(config, time.Now)
}

func NewWithClock(config Config, clock func() time.Time) *Application {
	memory := store.NewMemory()
	server := httpapi.New(httpapi.Deps{
		Apiary:     apiary.New(memory, clock),
		Inspection: inspection.New(memory, clock),
		Queen:      queen.New(memory, clock),
		Status:     status.New(memory, clock),
	})
	return &Application{Config: config, Store: memory, Server: server}
}

func (a *Application) Handler() http.Handler {
	return a.Server.Handler()
}
