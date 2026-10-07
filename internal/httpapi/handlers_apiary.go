package httpapi

import (
	"net/http"

	"example.com/apiary-colony-service/internal/apiary"
)

type registerApiaryRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	District string `json:"district"`
	Keeper   string `json:"keeper"`
}

func (s *Server) handleRegisterApiary(w http.ResponseWriter, r *http.Request) {
	var request registerApiaryRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, err)
		return
	}
	created, err := s.deps.Apiary.RegisterApiary(apiary.RegisterApiaryInput{
		ID:       request.ID,
		Name:     request.Name,
		District: request.District,
		Keeper:   request.Keeper,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleGetApiary(w http.ResponseWriter, r *http.Request) {
	found, err := s.deps.Apiary.Apiary(r.PathValue("apiaryId"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, found)
}
