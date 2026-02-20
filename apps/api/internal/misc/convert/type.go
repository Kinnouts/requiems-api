package convert

type Result struct {
	From    string  `json:"from"`
	To      string  `json:"to"`
	Input   float64 `json:"input"`
	Result  float64 `json:"result"`
	Formula string  `json:"formula"`
}

func (Result) IsData() {}

// UnitsResult maps each measurement category to its supported unit keys.
type UnitsResult struct {
	Length      []string `json:"length"`
	Weight      []string `json:"weight"`
	Volume      []string `json:"volume"`
	Temperature []string `json:"temperature"`
	Area        []string `json:"area"`
	Speed       []string `json:"speed"`
}

func (UnitsResult) IsData() {}
