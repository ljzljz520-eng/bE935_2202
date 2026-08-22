package workflow

import (
	"fmt"
	"lawdrive/internal/domain"
	"sort"
	"time"
)

type TimelineItem struct {
	At     time.Time
	Kind   string
	Actor  string
	FileID string
	Text   string
}

func (s *Service) Timeline(caseID string) ([]TimelineItem, error) {
	audits, e := s.repo.AllAudits()
	if e != nil {
		return nil, e
	}
	out := make([]TimelineItem, 0)
	for _, a := range audits {
		if a.CaseID == caseID {
			out = append(out, TimelineItem{At: a.At, Kind: a.Action, Actor: a.Actor, FileID: a.FileID, Text: a.Detail})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out, nil
}
func TimelineText(items []TimelineItem) string {
	parts := make([]string, 0, len(items))
	for _, it := range items {
		parts = append(parts, fmt.Sprintf("%s:%s", it.Kind, it.FileID))
	}
	return fmt.Sprint(parts)
}
func AddTimeline(items []TimelineItem, at time.Time, kind, actor, fileID, text string) []TimelineItem {
	return append(items, TimelineItem{At: at, Kind: kind, Actor: actor, FileID: fileID, Text: text})
}
func FilterTimeline(items []TimelineItem, from, to time.Time) []TimelineItem {
	out := make([]TimelineItem, 0)
	for _, it := range items {
		if !from.IsZero() && it.At.Before(from) {
			continue
		}
		if !to.IsZero() && it.At.After(to) {
			continue
		}
		out = append(out, it)
	}
	return out
}
func EventTimeline(e domain.Event, at time.Time) TimelineItem {
	return TimelineItem{At: at, Kind: e.Name, Actor: e.Actor, FileID: e.FileID, Text: e.Data}
}
