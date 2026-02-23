package phone

import (
	"github.com/nyaruka/phonenumbers"
)

// Service provides phone number validation logic.
type Service struct{}

// NewService creates a new phone validation Service.
func NewService() *Service { return &Service{} }

// Validate parses and validates a phone number, returning structured metadata.
// When the number cannot be parsed or is not valid, Valid is false and the
// optional fields (Country, Type, Formatted) are omitted.
func (s *Service) Validate(number string) ValidateResponse {
	num, err := phonenumbers.Parse(number, "")
	if err != nil || !phonenumbers.IsValidNumber(num) {
		return ValidateResponse{Number: number, Valid: false}
	}

	return ValidateResponse{
		Number:    number,
		Valid:     true,
		Country:   phonenumbers.GetRegionCodeForNumber(num),
		Type:      numberType(phonenumbers.GetNumberType(num)),
		Formatted: phonenumbers.Format(num, phonenumbers.INTERNATIONAL),
	}
}

// numberType converts a phonenumbers type constant to a human-readable string.
func numberType(t phonenumbers.PhoneNumberType) string {
	switch t {
	case phonenumbers.MOBILE:
		return "mobile"
	case phonenumbers.FIXED_LINE:
		return "landline"
	case phonenumbers.FIXED_LINE_OR_MOBILE:
		return "landline_or_mobile"
	case phonenumbers.TOLL_FREE:
		return "toll_free"
	case phonenumbers.PREMIUM_RATE:
		return "premium_rate"
	case phonenumbers.SHARED_COST:
		return "shared_cost"
	case phonenumbers.VOIP:
		return "voip"
	case phonenumbers.PERSONAL_NUMBER:
		return "personal_number"
	case phonenumbers.PAGER:
		return "pager"
	case phonenumbers.UAN:
		return "uan"
	case phonenumbers.VOICEMAIL:
		return "voicemail"
	default:
		return "unknown"
	}
}
