package base64 //nolint:revive // package name matches the service it tests; renaming would obscure intent

import (
	"errors"
	"net/http"
	"testing"

	"requiems-api/platform/httpx"
)

func TestService_Encode(t *testing.T) {
	svc := NewService()

	tests := []struct {
		name    string
		value   string
		variant string
		want    string
	}{
		{
			name:  "standard encoding",
			value: "Hello, world!",
			want:  "SGVsbG8sIHdvcmxkIQ==",
		},
		{
			name:    "standard encoding explicit",
			value:   "Hello, world!",
			variant: "standard",
			want:    "SGVsbG8sIHdvcmxkIQ==",
		},
		{
			name:    "url-safe encoding",
			value:   "Hello, world!",
			variant: "url",
			want:    "SGVsbG8sIHdvcmxkIQ==",
		},
		{
			name:    "url-safe encoding with url-unsafe characters",
			value:   "\xfb\xff\xfe",
			variant: "url",
			want:    "-__-",
		},
		{
			name:    "standard encoding with url-unsafe characters",
			value:   "\xfb\xff\xfe",
			variant: "standard",
			want:    "+//+",
		},
		{
			name:  "empty string",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.Encode(tt.value, tt.variant)
			if got.Result != tt.want {
				t.Errorf("Encode(%q, %q) = %q, want %q", tt.value, tt.variant, got.Result, tt.want)
			}
		})
	}
}

func TestService_Decode(t *testing.T) {
	svc := NewService()

	tests := []struct {
		name    string
		value   string
		variant string
		want    string
		wantErr bool
	}{
		{
			name:  "standard decoding",
			value: "SGVsbG8sIHdvcmxkIQ==",
			want:  "Hello, world!",
		},
		{
			name:    "standard decoding explicit",
			value:   "SGVsbG8sIHdvcmxkIQ==",
			variant: "standard",
			want:    "Hello, world!",
		},
		{
			name:    "url-safe decoding",
			value:   "SGVsbG8sIHdvcmxkIQ==",
			variant: "url",
			want:    "Hello, world!",
		},
		{
			name:    "url-safe encoding with url-unsafe characters",
			value:   "-__-",
			variant: "url",
			want:    "\xfb\xff\xfe",
		},
		{
			name:    "invalid standard base64",
			value:   "not-valid-base64!!!",
			wantErr: true,
		},
		{
			name:    "invalid url base64",
			value:   "not+valid+base64!!!",
			variant: "url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.Decode(tt.value, tt.variant)

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
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Result != tt.want {
				t.Errorf("Decode(%q, %q) = %q, want %q", tt.value, tt.variant, got.Result, tt.want)
			}
		})
	}
}
