package report

import (
	"sort"
	"strings"
)

type Row struct {
	Key   string
	Value int
}

func Rows(values map[string]int) []Row {
	out := make([]Row, 0, len(values))
	for key, value := range values {
		out = append(out, Row{Key: key, Value: value})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}
func Table(summary Summary) [][]string {
	rows := [][]string{{"dimension", "key", "count"}}
	for _, row := range Rows(summary.ByStore) {
		rows = append(rows, []string{"store", row.Key, itoa(row.Value)})
	}
	for _, row := range Rows(summary.BySeverity) {
		rows = append(rows, []string{"severity", row.Key, itoa(row.Value)})
	}
	return rows
}
func CSV(summary Summary) string {
	lines := make([]string, 0)
	for _, row := range Table(summary) {
		lines = append(lines, strings.Join(row, ","))
	}
	return strings.Join(lines, "\n") + "\n"
}
func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	digits := []byte{}
	for value > 0 {
		digits = append([]byte{byte('0' + value%10)}, digits...)
		value /= 10
	}
	return string(digits)
}
