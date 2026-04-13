package main

import "testing"

func TestNormaliseScheme(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"VISA", "visa"},
		{"visa", "visa"},
		{"MASTERCARD", "mastercard"},
		{"Master Card", "mastercard"},
		{"AMEX", "amex"},
		{"American Express", "amex"},
		{"DISCOVER", "discover"},
		{"JCB", "jcb"},
		{"DINERS CLUB", "diners"},
		{"UNIONPAY", "unionpay"},
		{"Union Pay", "unionpay"},
		{"MAESTRO", "maestro"},
		{"MIR", "mir"},
		{"RUPAY", "rupay"},
		{"PRIVATE LABEL", "private_label"},
		{"unknown", "unknown"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			if got := normaliseScheme(tt.input); got != tt.want {
				t.Fatalf("normaliseScheme(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormaliseCardType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"CREDIT", "credit"},
		{"Credit", "credit"},
		{"DEBIT", "debit"},
		{"PREPAID", "prepaid"},
		{"Pre-Paid", "prepaid"},
		{"CHARGE", "charge"},
		{"Charge Card", "charge"},
		{"other", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			if got := normaliseCardType(tt.input); got != tt.want {
				t.Fatalf("normaliseCardType(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormaliseCardLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"CLASSIC", "classic"},
		{"Normal", "classic"},
		{"GOLD", "gold"},
		{"Premier", "gold"},
		{"PLATINUM", "platinum"},
		{"World", "platinum"},
		{"INFINITE", "infinite"},
		{"BLACK", "infinite"},
		{"BUSINESS", "business"},
		{"CORPORATE", "corporate"},
		{"SIGNATURE", "signature"},
		{"STANDARD", "standard"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			if got := normaliseCardLevel(tt.input); got != tt.want {
				t.Fatalf("normaliseCardLevel(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestDetectScheme(t *testing.T) {
	t.Parallel()

	tests := []struct {
		bin  string
		want string
	}{
		{"411111", "visa"},       // Visa starts with 4
		{"412345", "visa"},
		{"510000", "mastercard"}, // Mastercard 51-55
		{"555555", "mastercard"},
		{"222100", "mastercard"}, // Mastercard 2-series
		{"340000", "amex"},       // Amex 34
		{"370000", "amex"},       // Amex 37
		{"352800", "jcb"},        // JCB 3528-3589
		{"300000", "diners"},     // Diners 3000-3059
		{"360000", "diners"},     // Diners 36
		{"601100", "discover"},   // 6011 prefix → Discover (takes precedence over RuPay 60)
		{"600000", "rupay"},      // RuPay 60 (not 6011)
		{"621000", "unionpay"},   // UnionPay 62
		{"220000", "mir"},        // Mir 2200-2204
		{"63", ""},               // too short
	}

	for _, tt := range tests {
		t.Run(tt.bin, func(t *testing.T) {
			t.Parallel()
			if got := detectScheme(tt.bin); got != tt.want {
				t.Fatalf("detectScheme(%q) = %q, want %q", tt.bin, got, tt.want)
			}
		})
	}
}

func TestNormalise_PrepaidSetsFlag(t *testing.T) {
	t.Parallel()

	r := RawBINRecord{
		BINPrefix:   "411111",
		Scheme:      "VISA",
		CardType:    "PREPAID",
		CountryCode: "US",
		Confidence:  0.75,
	}

	got := normalise(r)
	if !got.Prepaid {
		t.Fatal("expected Prepaid = true for card type PREPAID")
	}
	if got.Scheme != "visa" {
		t.Fatalf("Scheme = %q, want %q", got.Scheme, "visa")
	}
}

func TestNormalise_InvalidCountryCodeCleared(t *testing.T) {
	t.Parallel()

	r := RawBINRecord{
		BINPrefix:   "411111",
		CountryCode: "USA", // 3 letters — not valid 2-letter ISO
		CountryName: "United States",
		Confidence:  0.5,
	}

	got := normalise(r)
	if got.CountryCode != "" {
		t.Fatalf("CountryCode = %q, want empty for 3-letter code", got.CountryCode)
	}
	if got.CountryName != "" {
		t.Fatalf("CountryName = %q, want empty when CountryCode is cleared", got.CountryName)
	}
}

func TestMergeRecords_HigherConfidenceWins(t *testing.T) {
	t.Parallel()

	records := []RawBINRecord{
		{BINPrefix: "411111", IssuerName: "Low Conf Bank", Confidence: 0.5, Source: "src-a"},
		{BINPrefix: "411111", IssuerName: "High Conf Bank", Confidence: 0.9, Source: "src-b"},
	}

	merged := mergeRecords(records)
	r, ok := merged["411111"]
	if !ok {
		t.Fatal("expected merged record for BIN 411111")
	}
	if r.IssuerName != "High Conf Bank" {
		t.Fatalf("IssuerName = %q, want %q", r.IssuerName, "High Conf Bank")
	}
}

func TestMergeRecords_FillsEmptyFieldsFromLoser(t *testing.T) {
	t.Parallel()

	records := []RawBINRecord{
		// winner by confidence but missing IssuerURL
		{BINPrefix: "411111", IssuerName: "Winner Bank", IssuerURL: "", Confidence: 0.9, Source: "src-a"},
		// loser has IssuerURL
		{BINPrefix: "411111", IssuerName: "", IssuerURL: "https://loser.com", Confidence: 0.5, Source: "src-b"},
	}

	merged := mergeRecords(records)
	r := merged["411111"]
	if r.IssuerURL != "https://loser.com" {
		t.Fatalf("IssuerURL = %q, want filled from loser", r.IssuerURL)
	}
}

func TestCombineSources(t *testing.T) {
	t.Parallel()

	tests := []struct {
		a, b string
		want string
	}{
		{"src-a", "src-a", "src-a"},               // same → no duplicate
		{"src-a", "src-b", "src-a,src-b"},          // two distinct sources
		{"src-a,src-b", "src-b", "src-a,src-b"},   // already contains src-b
	}

	for _, tt := range tests {
		t.Run(tt.a+"/"+tt.b, func(t *testing.T) {
			t.Parallel()
			if got := combineSources(tt.a, tt.b); got != tt.want {
				t.Fatalf("combineSources(%q, %q) = %q, want %q", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
