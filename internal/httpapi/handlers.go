package httpapi

import (
	"coldrisk.local/console/internal/domain"
	"encoding/json"
	"io"
	"net/http"
)

func (s *Server) todo(w http.ResponseWriter, r *http.Request) {
	principal := r.Header.Get("X-Actor")
	storeID := r.URL.Query().Get("store")
	items, err := s.Service.Todo(principal, domain.RecordFilter{StoreID: storeID})
	if err != nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, items)
}
func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	principal := r.Header.Get("X-Actor")
	if r.Method == http.MethodGet {
		items, err := s.Service.Search(principal, domain.RecordFilter{Query: r.URL.Query().Get("q"), StoreID: r.URL.Query().Get("store"), IncludeArchived: true})
		if err != nil {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, items)
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var record domain.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	attachment := domain.NewAttachment("att-"+record.ID, record.ID, record.StoreID, "photo", "pending://"+record.ID, principal, record.CreatedAt)
	workflow := domain.NewWorkflow("wf-"+record.ID, record.ID, record.StoreID, principal, "")
	if err := s.Service.CreateRecord(principal, record, attachment, workflow, record.CreatedAt); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, record)
}
func (s *Server) importRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	principal := r.Header.Get("X-Actor")
	result, err := s.Service.Import(principal, r.Body, "import-time")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}
func (s *Server) report(w http.ResponseWriter, r *http.Request) {
	principal := r.Header.Get("X-Actor")
	summary, err := s.Service.BuildReport(principal, domain.RecordFilter{StoreID: r.URL.Query().Get("store"), IncludeArchived: true})
	if err != nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, summary)
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func decodeBody(body io.Reader, value any) error { return json.NewDecoder(body).Decode(value) }
