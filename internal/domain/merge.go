package domain

import "errors"

type Patch struct {
	Title       *string
	Description *string
	Severity    *string
	Assignee    *string
	Status      *string
}

func (r *Record) Apply(p Patch) error {
	if r == nil {
		return errors.New("record is nil")
	}
	if p.Title != nil {
		r.Title = *p.Title
	}
	if p.Description != nil {
		r.Description = *p.Description
	}
	if p.Severity != nil {
		r.Severity = NormalizeSeverity(*p.Severity)
	}
	if p.Assignee != nil {
		r.Assignee = *p.Assignee
	}
	if p.Status != nil {
		if err := r.Transition(*p.Status); err != nil {
			return err
		}
	}
	return r.Validate()
}
func (r Record) Diff(other Record) []string {
	out := []string{}
	if r.Title != other.Title {
		out = append(out, "title")
	}
	if r.Description != other.Description {
		out = append(out, "description")
	}
	if r.Severity != other.Severity {
		out = append(out, "severity")
	}
	if r.Status != other.Status {
		out = append(out, "status")
	}
	if r.Assignee != other.Assignee {
		out = append(out, "assignee")
	}
	return out
}
func MergePhotoRefs(left, right []string) []string {
	out := append([]string(nil), left...)
	seen := map[string]bool{}
	for _, id := range out {
		seen[id] = true
	}
	for _, id := range right {
		if id != "" && !seen[id] {
			out = append(out, id)
			seen[id] = true
		}
	}
	return out
}
func (r *Record) MergePhotos(items []string) error {
	if r == nil {
		return errors.New("record is nil")
	}
	r.PhotoRefs = MergePhotoRefs(r.PhotoRefs, items)
	return nil
}
func (r Record) CopyWithStatus(status string) (Record, error) {
	copy := CloneRecord(r)
	if err := copy.Transition(status); err != nil {
		return copy, err
	}
	return copy, nil
}
func (r Record) IsSameStore(other Record) bool { return r.StoreID == other.StoreID }
