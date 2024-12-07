package api

import "net/http"

func (s *Server) HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
