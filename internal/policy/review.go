package policy

import (
	"coldrisk.local/console/internal/domain"
	"errors"
	"sort"
)

type Decision struct {
	Allowed   bool
	Reason    string
	Principal string
	Action    string
	StoreID   string
}

func (p *Policy) Decide(principal, action, storeID string) Decision {
	allowed := p.Can(principal, action, storeID)
	reason := "denied"
	if allowed {
		reason = "allowed"
	} else if principal == "" {
		reason = "principal required"
	} else if storeID != "" && !p.HasStore(principal, storeID) {
		reason = "store outside grant"
	}
	return Decision{Allowed: allowed, Reason: reason, Principal: principal, Action: action, StoreID: storeID}
}
func (p *Policy) HasStore(principal, storeID string) bool {
	for _, id := range p.FilterStores(principal) {
		if id == storeID {
			return true
		}
	}
	return false
}
func (p *Policy) Require(principal, action, storeID string) error {
	decision := p.Decide(principal, action, storeID)
	if !decision.Allowed {
		return errors.New(decision.Reason)
	}
	return nil
}
func (p *Policy) CanAny(principal, action string, stores []string) bool {
	for _, storeID := range stores {
		if p.Can(principal, action, storeID) {
			return true
		}
	}
	return false
}
func (p *Policy) CanAll(principal, action string, stores []string) bool {
	for _, storeID := range stores {
		if !p.Can(principal, action, storeID) {
			return false
		}
	}
	return true
}
func (p *Policy) Intersection(left, right string) []string {
	a := p.FilterStores(left)
	b := p.FilterStores(right)
	set := map[string]bool{}
	for _, id := range b {
		set[id] = true
	}
	out := []string{}
	for _, id := range a {
		if set[id] {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}
func (p *Policy) Union(left, right string) []string {
	set := map[string]bool{}
	for _, id := range p.FilterStores(left) {
		set[id] = true
	}
	for _, id := range p.FilterStores(right) {
		set[id] = true
	}
	out := make([]string, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
func (p *Policy) Restrict(permission domain.Permission, stores []string) domain.Permission {
	allowed := map[string]bool{}
	for _, id := range stores {
		allowed[id] = true
	}
	original := permission.Stores
	permission.Stores = map[string]string{}
	for id, label := range original {
		if allowed[id] {
			permission.Stores[id] = label
		}
	}
	return permission
}
func (p *Policy) Reviewable(principal string, record domain.Record) bool {
	return record.IsActive() && domain.IsReviewable(record.Status) && p.Can(principal, ActionReview, record.StoreID)
}
func (p *Policy) CanModifyEvidence(principal string, record domain.Record) bool {
	return p.Can(principal, ActionReview, record.StoreID) && record.Status != domain.StatusArchived
}
func (p *Policy) CanClose(principal string, record domain.Record) bool {
	return p.Can(principal, ActionArchive, record.StoreID) && record.Status == domain.StatusPublished
}
func (p *Policy) PermissionCount() int {
	if p == nil {
		return 0
	}
	return len(p.permissions)
}
