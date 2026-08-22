package service

import (
	"coldrisk.local/console/internal/domain"
	"coldrisk.local/console/internal/policy"
	"coldrisk.local/console/internal/store"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkflowImportReport(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "risk.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	p := policy.New()
	p.Set(domain.Permission{Principal: "importer", Stores: map[string]string{"store-a": "import"}, Actions: map[string]bool{"import_records": true, "view_todo": true}})
	svc := New(st, p)
	good := domain.NewRecord("r1", "store-a", "Import one", "low", "t1")
	raw, _ := good.Marshal()
	input := strings.NewReader(string(raw) + "\n{bad}\n")
	result, err := svc.Import("importer", input, "t2")
	if err != nil {
		t.Fatal(err)
	}
	if result.Accepted != 1 || result.Rejected != 1 {
		t.Fatalf("unexpected import %#v", result)
	}
	summary, err := svc.BuildReport("importer", domain.RecordFilter{StoreID: "store-a", IncludeArchived: true})
	if err != nil || summary.Total != 1 {
		t.Fatalf("report failed %v %#v", err, summary)
	}
}
