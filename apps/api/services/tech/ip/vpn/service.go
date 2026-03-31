package vpn

import (
	"context"
	"net"

	"github.com/bobadilla-tech/go-ip-intelligence/v2/ipi"
)

type Service struct {
	c *ipi.Client
}

func NewService(databasePath, asnDatabasePath string) (*Service, error) {
	client, err := ipi.New(
		ipi.WithDatabasePath(databasePath),
		ipi.WithASNDatabasePath(asnDatabasePath),
	)
	if err != nil {
		return nil, err
	}

	return &Service{
		c: client,
	}, nil
}

func (s *Service) CheckIP(ctx context.Context, ip net.IP) (IPCheckResponse, error) {
	result, err := s.c.Check(ctx, ip)
	if err != nil {
		return IPCheckResponse{}, err
	}
	return IPCheckResponse{
		IP:         ip.String(),
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
