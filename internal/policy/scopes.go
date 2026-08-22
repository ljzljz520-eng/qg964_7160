package policy

import "coldrisk.local/console/internal/domain"

type Scope struct {
	Principal string
	Stores    []string
	Actions   []string
}

func (p *Policy) BuildScope(principal string) Scope {
	permission, _ := p.Permission(principal)
	actions := make([]string, 0)
	for action, allowed := range permission.Actions {
		if allowed {
			actions = append(actions, action)
		}
	}
	sortStrings(actions)
	return Scope{Principal: principal, Stores: p.FilterStores(principal), Actions: actions}
}
func (s Scope) HasStore(storeID string) bool {
	for _, id := range s.Stores {
		if id == storeID {
			return true
		}
	}
	return false
}
func (s Scope) HasAction(action string) bool {
	for _, id := range s.Actions {
		if id == action || id == "*" {
			return true
		}
	}
	return false
}
func (s Scope) Can(action, storeID string) bool {
	return s.HasAction(action) && (storeID == "" || s.HasStore(storeID))
}
func (p *Policy) StoresForRecord(principal string, records []domain.Record) []domain.Record {
	scope := p.BuildScope(principal)
	out := make([]domain.Record, 0)
	for _, r := range records {
		if scope.HasStore(r.StoreID) {
			out = append(out, r)
		}
	}
	domain.SortRecords(out)
	return out
}
func (p *Policy) CheckRecord(principal, action string, record domain.Record) bool {
	return p.Can(principal, action, record.StoreID)
}
func (p *Policy) AuthorizedActionSet(principal string) map[string]bool {
	item, _ := p.Permission(principal)
	out := map[string]bool{}
	for k, v := range item.Actions {
		if v {
			out[k] = true
		}
	}
	return out
}
func sortStrings(items []string) {
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && items[j] < items[j-1]; j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
}
