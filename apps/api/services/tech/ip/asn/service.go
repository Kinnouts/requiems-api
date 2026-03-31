package asn

import (
	"context"

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

func (s *Service) CheckASN(ctx context.Context, ip string) (ASNResponse, error) {
	info, err := s.c.CheckASNString(ctx, ip)
	if err != nil {
		return ASNResponse{}, err
	}
	return ASNResponse{
		IP:     info.IP,
		ASN:    info.ASN,
		Org:    info.Org,
		ISP:    info.ISP,
		Domain: info.Domain,
		Route:  info.Route,
		Type:   info.Type,
	}, nil
}
