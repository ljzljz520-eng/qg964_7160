package httpapi

import (
	"coldrisk.local/console/internal/domain"
	"net/url"
)

func filterFromQuery(values url.Values) domain.RecordFilter {
	return domain.RecordFilter{StoreID: values.Get("store"), Status: values.Get("status"), Severity: values.Get("severity"), Query: values.Get("q"), IncludeArchived: values.Get("archived") == "true"}
}
func actorFromHeaders(values map[string][]string) string {
	if items := values["X-Actor"]; len(items) > 0 {
		return items[0]
	}
	return ""
}
