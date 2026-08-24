package console

import (
	"context"
	"solidstate-battery-control/internal/domain"
	"testing"
)

func TestCanPublish(t *testing.T) {
	c := Console{}
	if c.CanPublish(context.Background(), domain.User{Role: "operator", Active: true}, domain.CellLot{State: domain.LotQualified}) == nil {
		t.Fatal("operator published")
	}
	if c.CanPublish(context.Background(), domain.User{Role: "reviewer", Active: true}, domain.CellLot{State: domain.LotDraft}) == nil {
		t.Fatal("draft published")
	}
}
