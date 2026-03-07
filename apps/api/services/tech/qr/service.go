package qr

import (
	qrcode "github.com/skip2/go-qrcode"
)

// Service generates QR codes.
type Service struct{}

// NewService returns a new Service instance.
func NewService() *Service {
	return &Service{}
}

// Generate returns the raw PNG bytes for a QR code encoding data at the
// given pixel size.
func (s *Service) Generate(data string, size int) ([]byte, error) {
	return qrcode.Encode(data, qrcode.Medium, size)
}
