package report

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (Builder) Format(summary Summary) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("total=%d active=%d archived=%d\n", summary.Total, summary.Active, summary.Archived))
	for _, key := range summary.Stores() {
		b.WriteString("store " + key + "=" + fmt.Sprint(summary.ByStore[key]) + "\n")
	}
	for _, key := range summary.Severities() {
		b.WriteString("severity " + key + "=" + fmt.Sprint(summary.BySeverity[key]) + "\n")
	}
	for _, key := range summary.Statuses() {
		b.WriteString("status " + key + "=" + fmt.Sprint(summary.ByStatus[key]) + "\n")
	}
	for _, key := range summary.Actions() {
		b.WriteString("action " + key + "=" + fmt.Sprint(summary.ByAction[key]) + "\n")
	}
	return b.String()
}
func (Builder) JSON(summary Summary) ([]byte, error) { return json.Marshal(summary) }
func (Builder) Lines(summary Summary) []string {
	text := Builder{}.Format(summary)
	return strings.Split(strings.TrimSuffix(text, "\n"), "\n")
}
func (Builder) Compare(left, right Summary) map[string]int {
	return map[string]int{"total_delta": right.Total - left.Total, "active_delta": right.Active - left.Active, "archived_delta": right.Archived - left.Archived}
}
func (Builder) Empty() Summary {
	return Summary{ByStore: map[string]int{}, BySeverity: map[string]int{}, ByStatus: map[string]int{}, ByAction: map[string]int{}}
}
