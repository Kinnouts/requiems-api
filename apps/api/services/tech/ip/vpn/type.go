package vpn

import "github.com/bobadilla-tech/go-ip-intelligence/ipi"

type IPCheckResponse struct {
	Ip         string          `json:"ip"`
	IsVPN      bool            `json:"is_vpn"`
	IsProxy    bool            `json:"is_proxy"`
	IsTor      bool            `json:"is_tor"`
	IsHosting  bool            `json:"is_hosting"`
	Score      int             `json:"score"`
	Threat     ipi.ThreatLevel `json:"threat"`
	FraudScore int             `json:"fraud_score"`
	AsnOrg     string          `json:"asn_org"`
}

func (IPCheckResponse) IsData() {}
