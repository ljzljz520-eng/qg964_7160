package report

import (
	"coldrisk.local/console/internal/domain"
	"strings"
	"testing"
)

func TestSummaryFormatting(t *testing.T) {
	r := domain.NewRecord("r1", "store-a", "Door", "high", "t")
	e := domain.NewAuditEvent("e1", "r1", "a", "created", "store-a", "x", "t")
	summary := NewBuilder().Build([]domain.Record{r}, []domain.AuditEvent{e})
	text := NewBuilder().Format(summary)
	if !strings.Contains(text, "total=1") || !strings.Contains(text, "store store-a=1") {
		t.Fatal(text)
	}
	if !strings.HasSuffix(CSV(summary), "\n") {
		t.Fatal("csv missing newline")
	}
}
