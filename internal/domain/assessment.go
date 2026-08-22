package domain

import (
	"fmt"
	"strings"
)

type Assessment struct {
	RecordID string            `json:"record_id"`
	Score    int               `json:"score"`
	Band     string            `json:"band"`
	Signals  map[string]string `json:"signals"`
}

func (r Record) Assess() Assessment {
	score := r.SeverityRank() * 20
	signals := map[string]string{"severity": r.Severity, "status": r.Status}
	if r.Assignee == "" {
		score += 10
		signals["assignee"] = "missing"
	} else {
		signals["assignee"] = "assigned"
	}
	if len(r.PhotoRefs) == 0 {
		score += 15
		signals["evidence"] = "missing"
	} else {
		signals["evidence"] = "present"
	}
	if r.Description == "" {
		score += 5
		signals["description"] = "missing"
	} else {
		signals["description"] = "present"
	}
	band := "normal"
	if score >= 70 {
		band = "critical"
	} else if score >= 45 {
		band = "elevated"
	}
	return Assessment{RecordID: r.ID, Score: score, Band: band, Signals: signals}
}

func (r Record) RiskKey() string {
	return strings.ToLower(strings.TrimSpace(r.StoreID)) + ":" + strings.ToLower(strings.TrimSpace(r.ID))
}
func (r Record) Brief() string { return fmt.Sprintf("%s %s %s", r.ID, r.Status, r.Severity) }
func (r Record) IsCritical() bool {
	return r.Severity == SeverityCritical || r.Assess().Band == "critical"
}
func (r Record) RequiresPhoto() bool {
	return r.Severity == SeverityHigh || r.Severity == SeverityCritical
}
func (r Record) ReadyForReview() bool {
	return r.Status == StatusOpen && r.Assignee != "" && len(r.PhotoRefs) > 0
}
func (r Record) ReadyForArchive() bool { return r.Status == StatusPublished && r.Description != "" }
func (r *Record) Normalize() {
	if r == nil {
		return
	}
	r.ID = strings.TrimSpace(r.ID)
	r.StoreID = strings.TrimSpace(r.StoreID)
	r.Title = strings.TrimSpace(r.Title)
	r.Description = strings.TrimSpace(r.Description)
	r.Severity = strings.ToLower(strings.TrimSpace(r.Severity))
	r.Status = strings.ToLower(strings.TrimSpace(r.Status))
	r.Assignee = strings.TrimSpace(r.Assignee)
}
func (r Record) ValidateBusiness() error {
	if err := r.Validate(); err != nil {
		return err
	}
	if r.RequiresPhoto() && len(r.PhotoRefs) == 0 {
		return fmt.Errorf("severity %s requires evidence", r.Severity)
	}
	return nil
}
func (a Assessment) Action() string {
	switch a.Band {
	case "critical":
		return "escalate"
	case "elevated":
		return "review"
	default:
		return "monitor"
	}
}
func (a Assessment) StableFields() []string {
	out := []string{}
	for _, key := range []string{"assignee", "description", "evidence", "severity", "status"} {
		if value, ok := a.Signals[key]; ok {
			out = append(out, key+"="+value)
		}
	}
	return out
}
func (a Assessment) Summary() string { return fmt.Sprintf("%s:%d:%s", a.RecordID, a.Score, a.Band) }
func AssessmentSort(items []Assessment) {
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && items[j].Score > items[j-1].Score; j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
}

type Timeline struct {
	RecordID string
	Events   []AuditEvent
}

func BuildTimeline(recordID string, events []AuditEvent) Timeline {
	filtered := FilterAuditEvents(events, recordID)
	return Timeline{RecordID: recordID, Events: filtered}
}
func (t Timeline) LastAction() string {
	if len(t.Events) == 0 {
		return ""
	}
	return t.Events[len(t.Events)-1].Action
}
func (t Timeline) HasAction(action string) bool {
	for _, event := range t.Events {
		if event.Action == action {
			return true
		}
	}
	return false
}
func (t Timeline) Actors() []string {
	out := []string{}
	seen := map[string]bool{}
	for _, event := range t.Events {
		if !seen[event.Actor] {
			out = append(out, event.Actor)
			seen[event.Actor] = true
		}
	}
	return out
}
func (t Timeline) CountAction(action string) int {
	count := 0
	for _, event := range t.Events {
		if event.Action == action {
			count++
		}
	}
	return count
}
