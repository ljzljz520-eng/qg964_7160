package domain

import "strings"

var StatusLabels = map[string]string{StatusOpen: "Open", StatusInReview: "In review", StatusConfirmed: "Confirmed", StatusPublished: "Published", StatusArchived: "Archived"}
var SeverityLabels = map[string]string{SeverityLow: "Low", SeverityMedium: "Medium", SeverityHigh: "High", SeverityCritical: "Critical"}

func StatusLabel(status string) string {
	if value, ok := StatusLabels[status]; ok {
		return value
	}
	return strings.Title(status)
}
func SeverityLabel(severity string) string {
	if value, ok := SeverityLabels[severity]; ok {
		return value
	}
	return strings.Title(severity)
}
func AllStatuses() []string {
	return []string{StatusOpen, StatusInReview, StatusConfirmed, StatusPublished, StatusArchived}
}
func AllSeverities() []string {
	return []string{SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical}
}
func IsTerminal(status string) bool    { return status == StatusArchived }
func IsReviewable(status string) bool  { return status == StatusOpen || status == StatusInReview }
func IsPublishable(status string) bool { return status == StatusConfirmed }
func IsArchivable(status string) bool  { return status == StatusPublished }
func NormalizeStatus(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if ValidStatus(value) {
		return value
	}
	return StatusOpen
}
func NormalizeSeverity(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if ValidSeverity(value) {
		return value
	}
	return SeverityMedium
}
