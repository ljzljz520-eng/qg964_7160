package service

import (
	"coldrisk.local/console/internal/domain"
	"coldrisk.local/console/internal/policy"
	"coldrisk.local/console/internal/store"
	"path/filepath"
	"testing"
)

func TestBusiness01Regression(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "risk.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	p := policy.New()
	p.Set(domain.Permission{Principal: "regional", Roles: []string{"regional_supervisor"}, Stores: map[string]string{"store-a": "repair"}, Actions: map[string]bool{"view_todo": true}})
	svc := New(st, p)
	for _, r := range []domain.Record{domain.NewRecord("r-a", "store-a", "Authorized repair", domain.SeverityHigh, "t"), domain.NewRecord("r-b", "store-b", "Other repair", domain.SeverityHigh, "t")} {
		if err := st.SaveRecord(r); err != nil {
			t.Fatal(err)
		}
	}
	items, err := svc.Todo("regional", domain.RecordFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].StoreID != "store-a" {
		t.Fatalf("regional scope leaked records: %#v", items)
	}
}
