package password

// Request holds the optional query parameters for the password endpoint.
// Defaults should be set before calling httpx.BindQuery.
type Request struct {
	Length    int  `query:"length"    validate:"min=8,max=128"`
	Uppercase bool `query:"uppercase"`
	Numbers   bool `query:"numbers"`
	Symbols   bool `query:"symbols"`
}

// Password is the response payload for the password generator.
type Password struct {
	Password string `json:"password"`
	Length   int    `json:"length"`
	Strength string `json:"strength"`
}

func (Password) IsData() {}
