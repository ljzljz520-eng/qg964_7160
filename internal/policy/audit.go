package policy

import (
	"coldrisk.local/console/internal/domain"
	"sort"
)

type AccessLog struct {
	Principal string
	Action    string
	StoreID   string
	Allowed   bool
	Reason    string
}

func (p *Policy) Audit(principal, action, storeID string) AccessLog {
	d := p.Decide(principal, action, storeID)
	return AccessLog{Principal: principal, Action: action, StoreID: storeID, Allowed: d.Allowed, Reason: d.Reason}
}
func (p *Policy) AuditBatch(principal, action string, stores []string) []AccessLog {
	out := make([]AccessLog, 0, len(stores))
	for _, storeID := range stores {
		out = append(out, p.Audit(principal, action, storeID))
	}
	return out
}
func LogsForStore(logs []AccessLog, storeID string) []AccessLog {
	out := []AccessLog{}
	for _, item := range logs {
		if item.StoreID == storeID {
			out = append(out, item)
		}
	}
	return out
}
func AllowedLogs(logs []AccessLog) []AccessLog {
	out := []AccessLog{}
	for _, item := range logs {
		if item.Allowed {
			out = append(out, item)
		}
	}
	return out
}
func SortAccessLogs(logs []AccessLog) {
	sort.Slice(logs, func(i, j int) bool {
		if logs[i].StoreID == logs[j].StoreID {
			return logs[i].Action < logs[j].Action
		}
		return logs[i].StoreID < logs[j].StoreID
	})
}
func (p *Policy) PermissionForRecord(principal string, record domain.Record) AccessLog {
	return p.Audit(principal, ActionReview, record.StoreID)
}
func AccessSummary(logs []AccessLog) map[string]int {
	out := map[string]int{"allowed": 0, "denied": 0}
	for _, item := range logs {
		if item.Allowed {
			out["allowed"]++
		} else {
			out["denied"]++
		}
	}
	return out
}
