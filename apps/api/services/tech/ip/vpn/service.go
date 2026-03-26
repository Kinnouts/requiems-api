package vpn

import (
	"context"

	"github.com/bobadilla-tech/go-ip-intelligence/ipi"
)

type Service struct {
	c *ipi.Client
}

func NewService() (*Service, error) {
	client, err := ipi.New(
		ipi.WithDatabasePath(""),
		ipi.WithASNDatabasePath(""),
	)
	if err != nil {
		return nil, err
	}

	return &Service{
		c: client,
	}, nil
}

func (s *Service) CheckIP(ctx context.Context, ip string) (IPCheckResponse, error) {
	result, err := s.c.CheckString(ctx, ip)
	if err != nil {
		return IPCheckResponse{}, err
	}
	return IPCheckResponse{
		Ip:         ip,
		IsVPN:      result.IsVPN,
		IsProxy:    result.IsProxy,
		IsTor:      result.IsTor,
		IsHosting:  result.IsHosting,
		Score:      result.Score,
		Threat:     result.Threat,
		FraudScore: result.FraudScore,
		AsnOrg:     result.AsnOrg,
	}, nil
}
