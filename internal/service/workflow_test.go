package service

import (
	"coldrisk.local/console/internal/domain"
	"coldrisk.local/console/internal/policy"
	"coldrisk.local/console/internal/store"
	"path/filepath"
	"testing"
)

func workflowService(t *testing.T) *Service {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "risk.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	p := policy.New()
	p.Set(domain.Permission{Principal: "supervisor", Roles: []string{"regional_supervisor"}, Stores: map[string]string{"store-a": "repair", "store-b": "repair"}, Actions: map[string]bool{"create_record": true, "review_record": true, "publish_record": true, "archive_record": true, "view_todo": true, "import_records": true}})
	return New(st, p)
}
func TestWorkflowCreateReviewArchive(t *testing.T) {
	svc := workflowService(t)
	record := domain.NewRecord("r1", "store-a", "Door alarm", domain.SeverityHigh, "t1")
	if err := svc.CreateRecord("supervisor", record, domain.NewAttachment("a1", "r1", "store-a", "photo", "file://a1", "supervisor", "t1"), domain.NewWorkflow("w1", "r1", "store-a", "supervisor", "t2"), "t1"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Review("supervisor", "r1", "t2"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Confirm("supervisor", "r1", "t3"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Publish("supervisor", "r1", "t4"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Archive("supervisor", "r1", "t5"); err != nil {
		t.Fatal(err)
	}
}
