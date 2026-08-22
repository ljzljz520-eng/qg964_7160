package store

import (
	"coldrisk.local/console/internal/domain"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "risk.db")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := domain.NewRecord("r1", "store-a", "Evaporator", "high", "t1")
	attachment := domain.NewAttachment("a1", "r1", "store-a", "photo", "file://r1", "alice", "t1")
	workflow := domain.NewWorkflow("w1", "r1", "store-a", "alice", "t2")
	audit := domain.NewAuditEvent("e1", "r1", "alice", "created", "store-a", "created", "t1")
	for _, err = range []error{st.SaveRecord(record), st.SaveAttachment(attachment), st.SaveWorkflow(workflow), st.SaveAudit(audit)} {
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	st, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if _, err = st.GetRecord("r1"); err != nil {
		t.Fatal(err)
	}
	items, _ := st.AttachmentsFor("r1")
	if len(items) != 1 {
		t.Fatal("attachment missing")
	}
	if _, err = st.WorkflowFor("r1"); err != nil {
		t.Fatal(err)
	}
	audits, _ := st.AuditFor("r1")
	if len(audits) != 1 {
		t.Fatal("audit missing")
	}
}
func TestStoreSearch(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "risk.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	r := domain.NewRecord("r1", "store-a", "Door", "medium", "t")
	r.Description = "temperature"
	if err := st.SaveRecord(r); err != nil {
		t.Fatal(err)
	}
	items, err := st.Search(domain.RecordFilter{Query: "temp"})
	if err != nil || len(items) != 1 {
		t.Fatalf("search failed %v %#v", err, items)
	}
}
