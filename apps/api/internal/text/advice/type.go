package advice

type Advice struct {
	ID   int    `json:"id"`
	Text string `json:"advice"`
}

func (Advice) IsData() {}
