package service

import (
	"coldrisk.local/console/internal/domain"
	"coldrisk.local/console/internal/policy"
	"coldrisk.local/console/internal/store"
	"path/filepath"
	"testing"
)

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "risk.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	p := policy.New()
	p.Set(domain.Permission{Principal: "reviewer", Stores: map[string]string{"store-a": "repair"}, Actions: map[string]bool{"view_todo": true, "review_record": true, "publish_record": true}})
	svc := New(st, p)
	r := domain.NewRecord("r1", "store-a", "Compressor", "medium", "t1")
	r.Description = "initial"
	if err := st.SaveRecord(r); err != nil {
		t.Fatal(err)
	}
	items, err := svc.Search("reviewer", domain.RecordFilter{StoreID: "store-a"})
	if err != nil || len(items) != 1 {
		t.Fatalf("search failed %v", err)
	}
	if _, err := svc.UpdateDescription("reviewer", "r1", "updated", "t2"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Publish("reviewer", "r1", "t3"); err == nil {
		t.Fatal("publish should require confirmed state")
	}
}
