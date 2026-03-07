package randomuser

import (
	"fmt"
	"math/rand/v2"
	"net/url"
	"strings"
)

var firstNames = []string{
	"Alice", "Bob", "Carol", "David", "Emma",
	"Frank", "Grace", "Henry", "Isabella", "James",
	"Karen", "Liam", "Mia", "Noah", "Olivia",
	"Peter", "Quinn", "Rachel", "Samuel", "Tara",
	"Uma", "Victor", "Wendy", "Xander", "Yara",
	"Zoe", "Aaron", "Bella", "Carlos", "Diana",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones",
	"Garcia", "Miller", "Davis", "Rodriguez", "Martinez",
	"Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson",
	"Thomas", "Taylor", "Moore", "Jackson", "Martin",
	"Lee", "Perez", "Thompson", "White", "Harris",
	"Sanchez", "Clark", "Lewis", "Robinson", "Walker",
}

var emailDomains = []string{
	"example.com", "mail.com", "test.org", "sample.net",
	"demo.io", "fake.com", "placeholder.dev",
}

var streetNames = []string{
	"Main", "Oak", "Maple", "Cedar", "Pine",
	"Elm", "Washington", "Lake", "Hill", "River",
	"Sunset", "Park", "Forest", "Meadow", "Spring",
	"Valley", "Summit", "Highland", "Harbor", "Ridge",
}

var streetTypes = []string{
	"Street", "Avenue", "Boulevard", "Lane", "Drive",
	"Court", "Place", "Road", "Way", "Circle",
}

var cities = []string{
	"Springfield", "Riverside", "Fairview", "Madison", "Georgetown",
	"Franklin", "Greenville", "Bristol", "Clinton", "Salem",
	"Centerville", "Lexington", "Ashland", "Burlington", "Manchester",
	"Newport", "Arlington", "Bloomington", "Clayton", "Dover",
}

// states maps US state names to their two-letter abbreviations.
var states = [][2]string{
	{"Alabama", "AL"}, {"Alaska", "AK"}, {"Arizona", "AZ"},
	{"California", "CA"}, {"Colorado", "CO"}, {"Florida", "FL"},
	{"Georgia", "GA"}, {"Illinois", "IL"}, {"Indiana", "IN"},
	{"Massachusetts", "MA"}, {"Michigan", "MI"}, {"Minnesota", "MN"},
	{"New York", "NY"}, {"North Carolina", "NC"}, {"Ohio", "OH"},
	{"Oregon", "OR"}, {"Pennsylvania", "PA"}, {"Texas", "TX"},
	{"Virginia", "VA"}, {"Washington", "WA"},
}

// Service generates random fake user data.
type Service struct{}

// NewService returns a new Service.
func NewService() *Service {
	return &Service{}
}

// Generate returns a randomly generated User.
func (s *Service) Generate() User {
	first := pick(firstNames)
	last := pick(lastNames)
	name := first + " " + last

	email := strings.ToLower(first+"."+last) + "@" + pick(emailDomains)

	state := states[rand.IntN(len(states))]
	address := Address{
		Street:  fmt.Sprintf("%d %s %s", 100+rand.IntN(9900), pick(streetNames), pick(streetTypes)),
		City:    pick(cities),
		State:   state[0],
		Zip:     fmt.Sprintf("%05d", 10000+rand.IntN(90000)),
		Country: "United States",
	}

	// Phone: US format +1-555-XXX-XXXX (555 is the conventional fake area code).
	phone := fmt.Sprintf("+1-555-%03d-%04d", rand.IntN(1000), rand.IntN(10000))

	avatar := "https://api.dicebear.com/9.x/identicon/svg?seed=" + url.QueryEscape(name)

	return User{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Address: address,
		Avatar:  avatar,
	}
}

func pick[T any](slice []T) T {
	return slice[rand.IntN(len(slice))]
}
