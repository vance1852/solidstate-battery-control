package domain

import "strings"

type LotFilter struct {
	Query         string
	States        []LotState
	Limit, Offset int
}

func (f LotFilter) Match(l CellLot) bool {
	if f.Query != "" && !strings.Contains(strings.ToLower(l.Code), strings.ToLower(f.Query)) {
		return false
	}
	if len(f.States) == 0 {
		return true
	}
	for _, s := range f.States {
		if l.State == s {
			return true
		}
	}
	return false
}
func SortLots(lots []CellLot, newest bool) []CellLot {
	out := append([]CellLot(nil), lots...)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			swap := out[i].Code > out[j].Code
			if newest {
				swap = out[i].UpdatedAt.Before(out[j].UpdatedAt)
			}
			if swap {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
func PageLots(lots []CellLot, limit, offset int) []CellLot {
	if limit <= 0 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}
	if offset >= len(lots) {
		return []CellLot{}
	}
	end := offset + limit
	if end > len(lots) {
		end = len(lots)
	}
	return lots[offset:end]
}
func ParseStates(values []string) []LotState {
	out := make([]LotState, 0, len(values))
	for _, v := range values {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "draft":
			out = append(out, LotDraft)
		case "qualified":
			out = append(out, LotQualified)
		case "released":
			out = append(out, LotReleased)
		case "hold":
			out = append(out, LotHold)
		}
	}
	return out
}
