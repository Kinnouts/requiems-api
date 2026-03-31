package vpn

import (
	"context"
	"net"
	"testing"
)

func TestService_CheckIP(t *testing.T) {
	svc, err := NewService("", "")
	if err != nil {
		t.Skipf("VPN service not available: %v", err)
	}

	tests := []struct {
		name    string
		ip      string
		wantIP  string
		wantErr bool
	}{
		{
			name:    "valid IPv4",
			ip:      "8.8.8.8",
			wantIP:  "8.8.8.8",
			wantErr: false,
		},
		{
			name:    "valid IPv6",
			ip:      "2001:4860:4860::8888",
			wantIP:  "2001:4860:4860::8888",
			wantErr: false,
		},
		{
			name:    "another IPv4",
			ip:      "1.1.1.1",
			wantIP:  "1.1.1.1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			if ip == nil {
				t.Fatalf("failed to parse IP: %s", tt.ip)
			}

			result, err := svc.CheckIP(context.Background(), ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result.IP != tt.wantIP {
				t.Errorf("CheckIP() IP = %v, want %v", result.IP, tt.wantIP)
			}
		})
	}
}

func TestService_CheckIP_ResponseFields(t *testing.T) {
	svc, err := NewService("", "")
	if err != nil {
		t.Skipf("VPN service not available: %v", err)
	}

	ip := net.ParseIP("8.8.8.8")
	result, err := svc.CheckIP(context.Background(), ip)
	if err != nil {
		t.Fatalf("CheckIP() unexpected error: %v", err)
	}

	if result.IP == "" {
		t.Error("expected non-empty IP")
	}

	if result.Score < 0 {
		t.Errorf("expected non-negative score, got %d", result.Score)
	}

	validThreats := map[string]bool{
		"none":     true,
		"low":      true,
		"medium":   true,
		"high":     true,
		"critical": true,
	}
	if !validThreats[result.Threat.String()] {
		t.Errorf("invalid threat level: %s", result.Threat)
	}

	if result.FraudScore < 0 || result.FraudScore > 100 {
		t.Errorf("fraud_score out of range: %d", result.FraudScore)
	}
}

func TestService_CheckIP_InvalidService(t *testing.T) {
	_, err := NewService("/nonexistent/path", "/another/nonexistent")
	if err == nil {
		t.Error("expected error for invalid database paths")
	}
}
