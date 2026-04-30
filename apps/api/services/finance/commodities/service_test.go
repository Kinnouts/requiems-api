package commodities

import (
	"context"
	"testing"

	"requiems-api/platform/httpx"
)

// stubGetter implements Getter for service-layer unit tests.
type stubGetter struct {
	result CommodityPrice
	err    error
}

func (s *stubGetter) Get(_ context.Context, slug string) (CommodityPrice, error) {
	if s.err != nil {
		return CommodityPrice{}, s.err
	}
	r := s.result
	r.Commodity = slug
	return r, nil
}

// ---- HistoricalPrice ----

func TestHistoricalPrice_FieldsPresent(t *testing.T) {
	h := HistoricalPrice{Period: "2023", Price: 1940.54}
	if h.Period != "2023" {
		t.Errorf("expected period 2023, got %q", h.Period)
	}
	if h.Price != 1940.54 {
		t.Errorf("expected price 1940.54, got %v", h.Price)
	}
}

// ---- CommodityPrice ----

func TestCommodityPrice_IsData(t *testing.T) {
	// IsData() must be callable — verifies the interface is satisfied.
	var c CommodityPrice
	c.IsData()
}

func TestCommodityPrice_FullResponse(t *testing.T) {
	cp := CommodityPrice{
		Commodity: "gold",
		Name:      "Gold",
		Price:     2386.33,
		Unit:      "oz",
		Currency:  "USD",
		Change24h: 23.01,
		Historical: []HistoricalPrice{
			{Period: "2023", Price: 1940.54},
			{Period: "2022", Price: 1800.12},
		},
	}

	if cp.Commodity != "gold" {
		t.Errorf("expected commodity gold, got %q", cp.Commodity)
	}
	if cp.Price != 2386.33 {
		t.Errorf("expected price 2386.33, got %v", cp.Price)
	}
	if len(cp.Historical) != 2 {
		t.Errorf("expected 2 historical entries, got %d", len(cp.Historical))
	}
}

// ---- Getter stub ----

func TestStubGetter_ReturnsCommodity(t *testing.T) {
	stub := &stubGetter{result: CommodityPrice{
		Name:  "Gold",
		Price: 2386.33,
		Unit:  "oz",
	}}

	result, err := stub.Get(context.Background(), "gold")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Commodity != "gold" {
		t.Errorf("expected commodity gold, got %q", result.Commodity)
	}
}

func TestStubGetter_PropagatesError(t *testing.T) {
	stub := &stubGetter{err: &httpx.AppError{
		Status:  404,
		Code:    "not_found",
		Message: "commodity not found",
	}}

	_, err := stub.Get(context.Background(), "unknown")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
