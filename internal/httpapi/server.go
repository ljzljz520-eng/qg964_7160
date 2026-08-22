package httpapi

import (
	"coldrisk.local/console/internal/service"
	"net/http"
)

type Server struct {
	Service *service.Service
	Mux     *http.ServeMux
}

func New(svc *service.Service) *Server {
	s := &Server{Service: svc, Mux: http.NewServeMux()}
	s.routes()
	return s
}
func (s *Server) routes() {
	s.Mux.HandleFunc("/health", s.health)
	s.Mux.HandleFunc("/todo", s.todo)
	s.Mux.HandleFunc("/records", s.records)
	s.Mux.HandleFunc("/import", s.importRecords)
	s.Mux.HandleFunc("/report", s.report)
}
func (s *Server) Handler() http.Handler { return s.Mux }
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": s.Service != nil && s.Service.Health()})
}
