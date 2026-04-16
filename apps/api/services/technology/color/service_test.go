package color //nolint:revive // package name matches the service it tests; renaming would obscure intent

import (
	"errors"
	"net/http"
	"testing"

	"requiems-api/platform/httpx"
)

func TestService_Convert(t *testing.T) {
	svc := NewService()

	tests := []struct {
		name        string
		from        string
		to          string
		value       string
		wantInput   string
		wantResult  string
		wantHex     string
		wantRGB     string
		wantHSL     string
		wantCMYK    string
		wantErr     bool
		wantErrCode string
	}{
		{
			name:       "hex to hsl",
			from:       "hex",
			to:         "hsl",
			value:      "#ff5733",
			wantInput:  "#ff5733",
			wantResult: "hsl(11, 100%, 60%)",
			wantHex:    "#ff5733",
			wantRGB:    "rgb(255, 87, 51)",
			wantHSL:    "hsl(11, 100%, 60%)",
			wantCMYK:   "cmyk(0%, 66%, 80%, 0%)",
		},
		{
			name:       "hex to rgb",
			from:       "hex",
			to:         "rgb",
			value:      "#ff5733",
			wantInput:  "#ff5733",
			wantResult: "rgb(255, 87, 51)",
			wantHex:    "#ff5733",
			wantRGB:    "rgb(255, 87, 51)",
			wantHSL:    "hsl(11, 100%, 60%)",
			wantCMYK:   "cmyk(0%, 66%, 80%, 0%)",
		},
		{
			name:       "hex to cmyk",
			from:       "hex",
			to:         "cmyk",
			value:      "#ff5733",
			wantInput:  "#ff5733",
			wantResult: "cmyk(0%, 66%, 80%, 0%)",
			wantHex:    "#ff5733",
			wantRGB:    "rgb(255, 87, 51)",
			wantHSL:    "hsl(11, 100%, 60%)",
			wantCMYK:   "cmyk(0%, 66%, 80%, 0%)",
		},
		{
			name:       "hex to hex (identity)",
			from:       "hex",
			to:         "hex",
			value:      "#ff5733",
			wantInput:  "#ff5733",
			wantResult: "#ff5733",
			wantHex:    "#ff5733",
			wantRGB:    "rgb(255, 87, 51)",
			wantHSL:    "hsl(11, 100%, 60%)",
			wantCMYK:   "cmyk(0%, 66%, 80%, 0%)",
		},
		{
			name:       "shorthand hex #rgb expands correctly",
			from:       "hex",
			to:         "hex",
			value:      "#f53",
			wantInput:  "#f53",
			wantResult: "#ff5533",
			wantHex:    "#ff5533",
			wantRGB:    "rgb(255, 85, 51)",
			wantHSL:    "hsl(10, 100%, 60%)",
			wantCMYK:   "cmyk(0%, 67%, 80%, 0%)",
		},
		{
			name:       "rgb to hex",
			from:       "rgb",
			to:         "hex",
			value:      "rgb(255, 87, 51)",
			wantInput:  "rgb(255, 87, 51)",
			wantResult: "#ff5733",
			wantHex:    "#ff5733",
			wantRGB:    "rgb(255, 87, 51)",
			wantHSL:    "hsl(11, 100%, 60%)",
			wantCMYK:   "cmyk(0%, 66%, 80%, 0%)",
		},
		{
			name:       "rgb without spaces",
			from:       "rgb",
			to:         "hex",
			value:      "rgb(255,87,51)",
			wantInput:  "rgb(255,87,51)",
			wantResult: "#ff5733",
			wantHex:    "#ff5733",
			wantRGB:    "rgb(255, 87, 51)",
			wantHSL:    "hsl(11, 100%, 60%)",
			wantCMYK:   "cmyk(0%, 66%, 80%, 0%)",
		},
		{
			name:       "black",
			from:       "hex",
			to:         "rgb",
			value:      "#000000",
			wantInput:  "#000000",
			wantResult: "rgb(0, 0, 0)",
			wantHex:    "#000000",
			wantRGB:    "rgb(0, 0, 0)",
			wantHSL:    "hsl(0, 0%, 0%)",
			wantCMYK:   "cmyk(0%, 0%, 0%, 100%)",
		},
		{
			name:       "white",
			from:       "hex",
			to:         "rgb",
			value:      "#ffffff",
			wantInput:  "#ffffff",
			wantResult: "rgb(255, 255, 255)",
			wantHex:    "#ffffff",
			wantRGB:    "rgb(255, 255, 255)",
			wantHSL:    "hsl(0, 0%, 100%)",
			wantCMYK:   "cmyk(0%, 0%, 0%, 0%)",
		},
		{
			name:       "hsl to hex",
			from:       "hsl",
			to:         "hex",
			value:      "hsl(120, 100%, 50%)",
			wantInput:  "hsl(120, 100%, 50%)",
			wantResult: "#00ff00",
			wantHex:    "#00ff00",
			wantRGB:    "rgb(0, 255, 0)",
			wantHSL:    "hsl(120, 100%, 50%)",
			wantCMYK:   "cmyk(100%, 0%, 100%, 0%)",
		},
		{
			name:       "cmyk to hex",
			from:       "cmyk",
			to:         "hex",
			value:      "cmyk(0%, 66%, 80%, 0%)",
			wantInput:  "cmyk(0%, 66%, 80%, 0%)",
			wantResult: "#ff5733",
			wantHex:    "#ff5733",
			wantRGB:    "rgb(255, 87, 51)",
			wantHSL:    "hsl(11, 100%, 60%)",
			wantCMYK:   "cmyk(0%, 66%, 80%, 0%)",
		},
		{
			name:        "invalid hex",
			from:        "hex",
			to:          "rgb",
			value:       "notahex",
			wantErr:     true,
			wantErrCode: "invalid_color",
		},
		{
			name:        "invalid rgb",
			from:        "rgb",
			to:          "hex",
			value:       "rgb(a, b, c)",
			wantErr:     true,
			wantErrCode: "invalid_color",
		},
		{
			name:        "invalid hsl",
			from:        "hsl",
			to:          "hex",
			value:       "hsl(a, b%, c%)",
			wantErr:     true,
			wantErrCode: "invalid_color",
		},
		{
			name:        "invalid cmyk",
			from:        "cmyk",
			to:          "hex",
			value:       "cmyk(a, b, c, d)",
			wantErr:     true,
			wantErrCode: "invalid_color",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.Convert(tt.from, tt.to, tt.value)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var appErr *httpx.AppError
				if !errors.As(err, &appErr) {
					t.Fatalf("expected *httpx.AppError, got %T", err)
				}
				if appErr.Status != http.StatusUnprocessableEntity {
					t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, appErr.Status)
				}
				if tt.wantErrCode != "" && appErr.Code != tt.wantErrCode {
					t.Errorf("expected code %q, got %q", tt.wantErrCode, appErr.Code)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Input != tt.wantInput {
				t.Errorf("Input = %q, want %q", got.Input, tt.wantInput)
			}
			if got.Result != tt.wantResult {
				t.Errorf("Result = %q, want %q", got.Result, tt.wantResult)
			}
			if got.Formats.Hex != tt.wantHex {
				t.Errorf("Formats.Hex = %q, want %q", got.Formats.Hex, tt.wantHex)
			}
			if got.Formats.RGB != tt.wantRGB {
				t.Errorf("Formats.RGB = %q, want %q", got.Formats.RGB, tt.wantRGB)
			}
			if got.Formats.HSL != tt.wantHSL {
				t.Errorf("Formats.HSL = %q, want %q", got.Formats.HSL, tt.wantHSL)
			}
			if got.Formats.CMYK != tt.wantCMYK {
				t.Errorf("Formats.CMYK = %q, want %q", got.Formats.CMYK, tt.wantCMYK)
			}
		})
	}
}
