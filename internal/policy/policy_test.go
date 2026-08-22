package policy

import (
	"coldrisk.local/console/internal/domain"
	"testing"
)

func TestPolicyScopes(t *testing.T) {
	p := NewWithPermissions([]domain.Permission{{Principal: "alice", Roles: []string{"regional_supervisor"}, Stores: map[string]string{"store-a": "repair"}, Actions: map[string]bool{ActionViewTodo: true, ActionReview: true}}})
	if !p.Can("alice", ActionViewTodo, "store-a") {
		t.Fatal("expected allowed")
	}
	if p.Can("alice", ActionViewTodo, "store-b") {
		t.Fatal("unexpected allowed")
	}
	scope := p.BuildScope("alice")
	if !scope.HasStore("store-a") || scope.HasStore("store-b") {
		t.Fatal("bad scope")
	}
}
