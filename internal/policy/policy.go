package policy

import (
	"coldrisk.local/console/internal/domain"
	"sort"
	"strings"
)

const (
	ActionViewTodo = "view_todo"
	ActionCreate   = "create_record"
	ActionReview   = "review_record"
	ActionPublish  = "publish_record"
	ActionArchive  = "archive_record"
	ActionImport   = "import_records"
)

type Policy struct{ permissions map[string]domain.Permission }

func New() *Policy { return &Policy{permissions: map[string]domain.Permission{}} }
func NewWithPermissions(items []domain.Permission) *Policy {
	p := New()
	for _, item := range items {
		p.Set(item)
	}
	return p
}
func (p *Policy) Set(permission domain.Permission) {
	if p == nil {
		return
	}
	if permission.Stores == nil {
		permission.Stores = map[string]string{}
	}
	if permission.Actions == nil {
		permission.Actions = map[string]bool{}
	}
	p.permissions[permission.Principal] = permission
}
func (p *Policy) Remove(principal string) {
	if p != nil {
		delete(p.permissions, principal)
	}
}
func (p *Policy) Permission(principal string) (domain.Permission, bool) {
	if p == nil {
		return domain.Permission{}, false
	}
	item, ok := p.permissions[principal]
	return item, ok
}
func (p *Policy) Can(principal, action, storeID string) bool {
	permission, ok := p.Permission(principal)
	if !ok {
		return false
	}
	if !permission.Actions[action] && !permission.Actions["*"] {
		return false
	}
	if storeID == "" {
		return action == ActionImport
	}
	_, ok = permission.Stores[storeID]
	return ok
}
func (p *Policy) FilterStores(principal string) []string {
	permission, ok := p.Permission(principal)
	if !ok {
		return nil
	}
	stores := make([]string, 0, len(permission.Stores))
	for id := range permission.Stores {
		stores = append(stores, id)
	}
	sort.Strings(stores)
	return stores
}
func (p *Policy) Scope(principal string) domain.Permission {
	item, _ := p.Permission(principal)
	return item
}
func (p *Policy) Role(principal string) string {
	item, ok := p.Permission(principal)
	if !ok {
		return ""
	}
	roles := append([]string(nil), item.Roles...)
	sort.Strings(roles)
	if len(roles) == 0 {
		return ""
	}
	return roles[0]
}
func (p *Policy) Explain(principal, action, storeID string) string {
	if p.Can(principal, action, storeID) {
		return "allowed"
	}
	return "denied"
}
func (p *Policy) Principals() []string {
	out := make([]string, 0)
	if p == nil {
		return out
	}
	for id := range p.permissions {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
func (p *Policy) Grant(principal, storeID, label string) {
	item, _ := p.Permission(principal)
	if item.Principal == "" {
		item.Principal = principal
	}
	if item.Stores == nil {
		item.Stores = map[string]string{}
	}
	item.Stores[storeID] = label
	p.Set(item)
}
func (p *Policy) Allow(principal, action string) {
	item, _ := p.Permission(principal)
	if item.Principal == "" {
		item.Principal = principal
	}
	if item.Actions == nil {
		item.Actions = map[string]bool{}
	}
	item.Actions[action] = true
	p.Set(item)
}
func (p *Policy) IsRegionalSupervisor(principal string) bool {
	role := strings.ToLower(p.Role(principal))
	return role == "regional_supervisor" || role == "supervisor"
}
