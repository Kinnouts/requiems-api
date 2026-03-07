package randomuser

import (
	"net/url"

	faker "github.com/jaswdr/faker/v2"
)

// Service generates random fake user data.
type Service struct{}

// NewService returns a new Service.
func NewService() *Service {
	return &Service{}
}

// Generate returns a randomly generated User.
func (s *Service) Generate() User {
	f := faker.New()

	name := f.Person().Name()
	email := f.Internet().SafeEmail()
	phone := f.Phone().Number()

	address := Address{
		Street:  f.Address().StreetAddress(),
		City:    f.Address().City(),
		State:   f.Address().State(),
		Zip:     f.Address().PostCode(),
		Country: f.Address().Country(),
	}

	avatar := "https://api.dicebear.com/9.x/identicon/svg?seed=" + url.QueryEscape(name)

	return User{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Address: address,
		Avatar:  avatar,
	}
}
