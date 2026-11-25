package pricing

import "testing"

func TestMimoV26RatesFollowTheVendorPage(t *testing.T) {
	cases := []struct {
		model          string
		hit, miss, out float64
		v25Twin        string
	}{
		{"mimo-v2.6-pro", 0.025, 3, 6, "mimo-v2.5-pro"},
		{"mimo-v2.6-flash", 0.02, 1, 2, "mimo-v2.5"},
		{"mimo-v2.6-pro-ultraspeed", 0.25, 30, 60, ""},
	}
	for _, tc := range cases {
		got := CurrentRate("mimo", tc.model, "CNY")
		if got == nil || got.CacheHit != tc.hit || got.Input != tc.miss || got.Output != tc.out {
			t.Fatalf("%s rate = %+v, want %v/%v/%v", tc.model, got, tc.hit, tc.miss, tc.out)
		}
		if tc.v25Twin == "" {
			continue
		}
		old := CurrentRate("mimo", tc.v25Twin, "CNY")
		if old == nil || *old != *got {
			t.Fatalf("%s must keep its published rate beside %s: %+v vs %+v", tc.v25Twin, tc.model, old, got)
		}
	}
}

func TestHistoricalMiMoCatalogRatesKeepTheirFingerprints(t *testing.T) {
	cases := []struct {
		model string
		mode  string
		rate  RateCard
	}{
		{"mimo-v2.5-pro", BillingModePAYG, RateCard{CacheHit: 0.025, Input: 3, Output: 6, Currency: "CNY"}},
		{"mimo-v2.5", BillingModePAYG, RateCard{CacheHit: 0.02, Input: 1, Output: 2, Currency: "CNY"}},
		{"mimo-v2-flash", BillingModePAYG, RateCard{CacheHit: 0.07, Input: 0.70, Output: 2.10, Currency: "CNY"}},
		{"mimo-v2.5-pro", BillingModeSubscriptionEquivalent, RateCard{CacheHit: 0.025, Input: 3, Output: 6, Currency: "CNY"}},
		{"mimo-v2.5", BillingModeSubscriptionEquivalent, RateCard{CacheHit: 0.02, Input: 1, Output: 2, Currency: "CNY"}},
	}
	for _, tc := range cases {
		row, ok := LookupCatalog("mimo", tc.model, "CNY", tc.mode)
		if !ok || !row.Historical || row.BillingMode != tc.mode || row.Fingerprint != PricingFingerprint(tc.rate) {
			t.Errorf("%s (%s) lost its historical catalog rate: %+v", tc.model, tc.mode, row)
		}
	}
	row, ok := LookupCatalog("mimo", "mimo-v2.6-pro", "CNY", BillingModePAYG)
	if !ok || row.Historical {
		t.Errorf("MiMo V2.6 Pro must be current: %+v", row)
	}
}
