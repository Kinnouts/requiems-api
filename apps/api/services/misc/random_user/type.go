package randomuser

// Address holds the user's postal address fields.
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}

// User is the response model for the random-user endpoint.
type User struct {
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Phone   string  `json:"phone"`
	Address Address `json:"address"`
	Avatar  string  `json:"avatar"`
}

func (User) IsData() {}
