package httpapi

import (
	"net/http"

	"example.com/apiary-colony-service/internal/domain"
	"example.com/apiary-colony-service/internal/queen"
)

type replaceQueenRequest struct {
	Origin       string `json:"origin"`
	IntroducedAt string `json:"introducedAt"`
	Reason       string `json:"reason"`
}

func (s *Server) handleReplaceQueen(w http.ResponseWriter, r *http.Request) {
	var request replaceQueenRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, err)
		return
	}
	introducedAt, err := parseOptionalTime("introducedAt", request.IntroducedAt)
	if err != nil {
		writeError(w, err)
		return
	}
	input := queen.ReplaceInput{
		Origin: request.Origin,
		Reason: request.Reason,
	}
	if introducedAt != nil {
		input.IntroducedAt = *introducedAt
	}
	record, err := s.deps.Queen.Replace(r.PathValue("hiveId"), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (s *Server) handleListQueens(w http.ResponseWriter, r *http.Request) {
	items, err := s.deps.Queen.List(r.PathValue("hiveId"))
	if err != nil {
		writeError(w, err)
		return
	}
	if items == nil {
		items = []domain.Queen{}
	}
	writeJSON(w, http.StatusOK, items)
}
