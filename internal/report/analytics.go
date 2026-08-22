package report

import (
	"coldrisk.local/console/internal/domain"
	"sort"
)

type Ranking struct {
	Key   string
	Count int
	Score int
}

func RankSeverities(records []domain.Record) []Ranking {
	values := map[string]Ranking{}
	for _, record := range records {
		item := values[record.Severity]
		item.Key = record.Severity
		item.Count++
		item.Score += record.SeverityRank()
		values[record.Severity] = item
	}
	out := make([]Ranking, 0, len(values))
	for _, item := range values {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].Key < out[j].Key
		}
		return out[i].Score > out[j].Score
	})
	return out
}
func RankStores(records []domain.Record) []Ranking {
	values := map[string]Ranking{}
	for _, record := range records {
		item := values[record.StoreID]
		item.Key = record.StoreID
		item.Count++
		item.Score += record.SeverityRank()
		values[record.StoreID] = item
	}
	out := make([]Ranking, 0, len(values))
	for _, item := range values {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].Key < out[j].Key
		}
		return out[i].Score > out[j].Score
	})
	return out
}
func OpenCount(records []domain.Record) int {
	count := 0
	for _, record := range records {
		if record.Status == domain.StatusOpen {
			count++
		}
	}
	return count
}
func ReviewCount(records []domain.Record) int {
	count := 0
	for _, record := range records {
		if record.Status == domain.StatusInReview {
			count++
		}
	}
	return count
}
func ConfirmationRate(records []domain.Record) float64 {
	if len(records) == 0 {
		return 0
	}
	count := 0
	for _, record := range records {
		if record.Status == domain.StatusConfirmed || record.Status == domain.StatusPublished || record.Status == domain.StatusArchived {
			count++
		}
	}
	return float64(count) / float64(len(records))
}
func EvidenceRate(records []domain.Record) float64 {
	if len(records) == 0 {
		return 0
	}
	count := 0
	for _, record := range records {
		if len(record.PhotoRefs) > 0 {
			count++
		}
	}
	return float64(count) / float64(len(records))
}
func CriticalRecords(records []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, record := range records {
		if record.IsCritical() {
			out = append(out, record)
		}
	}
	domain.SortRecords(out)
	return out
}
func Unassigned(records []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, record := range records {
		if record.Assignee == "" {
			out = append(out, record)
		}
	}
	domain.SortRecords(out)
	return out
}
func SummaryForStore(records []domain.Record, storeID string) Summary {
	filtered := []domain.Record{}
	for _, record := range records {
		if record.StoreID == storeID {
			filtered = append(filtered, record)
		}
	}
	return NewBuilder().Build(filtered, nil)
}
