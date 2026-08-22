package service

import (
	"bufio"
	"coldrisk.local/console/internal/domain"
	"coldrisk.local/console/internal/policy"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type ImportResult struct {
	Accepted int      `json:"accepted"`
	Rejected int      `json:"rejected"`
	Errors   []string `json:"errors"`
}

func (s *Service) Import(principal string, input io.Reader, at string) (ImportResult, error) {
	result := ImportResult{Errors: []string{}}
	if err := s.ready(); err != nil {
		return result, err
	}
	if !s.Policy.Can(principal, policy.ActionImport, "") {
		return result, fmt.Errorf("import forbidden")
	}
	scanner := bufio.NewScanner(input)
	line := 0
	for scanner.Scan() {
		line++
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" {
			continue
		}
		var record domain.Record
		if err := json.Unmarshal([]byte(raw), &record); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("line %d: %v", line, err))
			continue
		}
		if err := record.Validate(); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("line %d: %v", line, err))
			continue
		}
		if err := s.Store.SaveRecord(record); err != nil {
			return result, err
		}
		workflow := domain.NewWorkflow("wf-"+record.ID, record.ID, record.StoreID, principal, "")
		if err := s.Store.SaveWorkflow(workflow); err != nil {
			return result, err
		}
		audit := domain.NewAuditEvent(fmt.Sprintf("audit-%s-import", record.ID), record.ID, principal, "imported", record.StoreID, "record imported", at)
		if err := s.Store.SaveAudit(audit); err != nil {
			return result, err
		}
		result.Accepted++
	}
	if err := scanner.Err(); err != nil {
		return result, err
	}
	return result, nil
}
func (s *Service) ImportRecords(principal string, records []domain.Record, at string) (ImportResult, error) {
	var b strings.Builder
	for _, r := range records {
		data, _ := json.Marshal(r)
		b.Write(data)
		b.WriteByte('\n')
	}
	return s.Import(principal, strings.NewReader(b.String()), at)
}
func (r ImportResult) Successful() bool  { return r.Rejected == 0 }
func (r ImportResult) ErrorText() string { return strings.Join(r.Errors, "; ") }
