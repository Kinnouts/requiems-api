package info

import (
	"context"
	"net"

	"github.com/bobadilla-tech/go-ip-intelligence/v2/ipi"
)

type Service struct {
	c *ipi.Client
}

func NewService(c *ipi.Client) *Service {
	if c == nil {
		return nil
	}
	return &Service{c: c}
}

func (s *Service) CheckInfo(ctx context.Context, ip string) (InfoResponse, error) {
	result, err := s.c.CheckString(ctx, ip)
	if err != nil {
		return InfoResponse{}, err
	}
	return InfoResponse{
		IP:          net.IP(result.IP).String(),
		Country:     result.Country,
		CountryCode: result.CountryCode,
		City:        result.City,
		ISP:         result.AsnOrg,
		IsVPN:       result.IsVPN,
	}, nil
}
