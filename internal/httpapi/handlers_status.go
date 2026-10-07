package httpapi

import (
	"net/http"

	"example.com/apiary-colony-service/internal/domain"
)

func (s *Server) handleApiaryStatus(w http.ResponseWriter, r *http.Request) {
	asOf, err := parseOptionalTime("asOf", r.URL.Query().Get("asOf"))
	if err != nil {
		writeError(w, err)
		return
	}
	report, err := s.deps.Status.ApiaryStatus(r.PathValue("apiaryId"), asOf)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, domain.NotFound("路径", r.URL.Path))
}
