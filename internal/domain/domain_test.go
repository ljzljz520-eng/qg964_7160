package domain

import "testing"

func TestRecordValidationAndTransitions(t *testing.T) {
	r := NewRecord("r1", "store-a", "Door alarm", "high", "2026-01-01T00:00:00Z")
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := r.Transition(StatusInReview); err != nil {
		t.Fatal(err)
	}
	if err := r.Transition(StatusConfirmed); err != nil {
		t.Fatal(err)
	}
	if err := r.Transition(StatusPublished); err != nil {
		t.Fatal(err)
	}
	if err := r.Transition(StatusArchived); err != nil {
		t.Fatal(err)
	}
	if r.IsActive() {
		t.Fatal("archived record is active")
	}
}
func TestRecordFilter(t *testing.T) {
	items := []Record{NewRecord("r2", "store-b", "Fan", "low", "t"), NewRecord("r1", "store-a", "Door", "critical", "t")}
	items[0].Description = "fan issue"
	got := FilterRecords(items, RecordFilter{StoreID: "store-a"})
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("unexpected filter %#v", got)
	}
}
