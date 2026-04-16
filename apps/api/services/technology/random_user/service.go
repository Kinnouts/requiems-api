package randomuser

import (
	"net/url"

	faker "github.com/jaswdr/faker/v2"
)

// Generates random fake user data.
type Service struct{}

func NewService() *Service {
	return &Service{}
}

// Returns a randomly generated User.
func (s *Service) Generate() User {
	f := faker.New()

	name := f.Person().Name()
	email := f.Internet().SafeEmail()
	phone := f.Phone().Number()

	fakerAddress := f.Address()

	address := Address{
		Street:  fakerAddress.StreetAddress(),
		City:    fakerAddress.City(),
		State:   fakerAddress.State(),
		Zip:     fakerAddress.PostCode(),
		Country: fakerAddress.Country(),
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
