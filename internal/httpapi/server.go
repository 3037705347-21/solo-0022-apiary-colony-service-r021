package httpapi

import (
	"net/http"

	"example.com/apiary-colony-service/internal/apiary"
	"example.com/apiary-colony-service/internal/inspection"
	"example.com/apiary-colony-service/internal/queen"
	"example.com/apiary-colony-service/internal/status"
)

type Deps struct {
	Apiary     *apiary.Service
	Inspection *inspection.Service
	Queen      *queen.Service
	Status     *status.Service
}

// Server 暴露蜂场协作服务，并把领域错误映射为稳定响应。
type Server struct {
	deps Deps
	mux  *http.ServeMux
}

func New(deps Deps) *Server {
	server := &Server{deps: deps, mux: http.NewServeMux()}
	server.register()
	return server
}

func (s *Server) Handler() http.Handler {
	return s.mux
}
