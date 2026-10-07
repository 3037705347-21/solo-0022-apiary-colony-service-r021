package httpapi

func (s *Server) register() {
	s.mux.HandleFunc("GET /healthz", s.handleHealth)

	s.mux.HandleFunc("POST /apiaries", s.handleRegisterApiary)
	s.mux.HandleFunc("GET /apiaries/{apiaryId}", s.handleGetApiary)
	s.mux.HandleFunc("GET /apiaries/{apiaryId}/status", s.handleApiaryStatus)
	s.mux.HandleFunc("POST /apiaries/{apiaryId}/hives", s.handleInstallHive)

	s.mux.HandleFunc("GET /hives/{hiveId}", s.handleGetHive)
	s.mux.HandleFunc("POST /hives/{hiveId}/state", s.handleAdvanceHiveState)
	s.mux.HandleFunc("POST /hives/{hiveId}/inspections", s.handleRecordInspection)
	s.mux.HandleFunc("GET /hives/{hiveId}/inspections", s.handleListInspections)
	s.mux.HandleFunc("POST /hives/{hiveId}/queens", s.handleReplaceQueen)
	s.mux.HandleFunc("GET /hives/{hiveId}/queens", s.handleListQueens)

	s.mux.HandleFunc("/", s.handleNotFound)
}
