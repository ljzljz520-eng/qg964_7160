package httpapi

import (
	"coldrisk.local/console/internal/domain"
	"coldrisk.local/console/internal/policy"
	"coldrisk.local/console/internal/service"
	"coldrisk.local/console/internal/store"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestHTTPTodoScope(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "risk.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	p := policy.New()
	p.Set(domain.Permission{Principal: "regional", Stores: map[string]string{"store-a": "repair"}, Actions: map[string]bool{"view_todo": true}})
	svc := service.New(st, p)
	if err := st.SaveRecord(domain.NewRecord("r1", "store-a", "Door", "high", "t")); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/todo", nil)
	req.Header.Set("X-Actor", "regional")
	res := httptest.NewRecorder()
	New(svc).Handler().ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatalf("status %d", res.Code)
	}
}
