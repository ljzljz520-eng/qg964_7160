package report

import (
	"coldrisk.local/console/internal/domain"
	"sort"
)

type Summary struct {
	Total      int            `json:"total"`
	Active     int            `json:"active"`
	Archived   int            `json:"archived"`
	ByStore    map[string]int `json:"by_store"`
	BySeverity map[string]int `json:"by_severity"`
	ByStatus   map[string]int `json:"by_status"`
	ByAction   map[string]int `json:"by_action"`
}
type Builder struct{}

func NewBuilder() Builder { return Builder{} }
func (Builder) Build(records []domain.Record, audits []domain.AuditEvent) Summary {
	out := Summary{ByStore: map[string]int{}, BySeverity: map[string]int{}, ByStatus: map[string]int{}, ByAction: map[string]int{}}
	for _, r := range records {
		out.Total++
		if r.Status == domain.StatusArchived {
			out.Archived++
		} else {
			out.Active++
		}
		out.ByStore[r.StoreID]++
		out.BySeverity[r.Severity]++
		out.ByStatus[r.Status]++
	}
	for _, a := range audits {
		out.ByAction[a.Action]++
	}
	return out
}
func (s Summary) Stores() []string     { return sortedKeys(s.ByStore) }
func (s Summary) Severities() []string { return sortedKeys(s.BySeverity) }
func (s Summary) Statuses() []string   { return sortedKeys(s.ByStatus) }
func (s Summary) Actions() []string    { return sortedKeys(s.ByAction) }
func sortedKeys(values map[string]int) []string {
	out := make([]string, 0, len(values))
	for key := range values {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}
