package httpapi

import (
	"net/http"

	"example.com/apiary-colony-service/internal/apiary"
)

type installHiveRequest struct {
	ID    string `json:"id"`
	Model string `json:"model"`
}

type advanceHiveStateRequest struct {
	Target string `json:"target"`
}

func (s *Server) handleInstallHive(w http.ResponseWriter, r *http.Request) {
	var request installHiveRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, err)
		return
	}
	created, err := s.deps.Apiary.InstallHive(r.PathValue("apiaryId"), apiary.InstallHiveInput{
		ID:    request.ID,
		Model: request.Model,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleGetHive(w http.ResponseWriter, r *http.Request) {
	found, err := s.deps.Apiary.Hive(r.PathValue("hiveId"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, found)
}

func (s *Server) handleAdvanceHiveState(w http.ResponseWriter, r *http.Request) {
	var request advanceHiveStateRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, err)
		return
	}
	transition, err := s.deps.Apiary.AdvanceHiveState(r.PathValue("hiveId"), request.Target)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, transition)
}
