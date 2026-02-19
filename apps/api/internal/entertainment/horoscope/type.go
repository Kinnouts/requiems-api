package horoscope

type Horoscope struct {
	Sign        string `json:"sign"`
	Date        string `json:"date"`
	Horoscope   string `json:"horoscope"`
	LuckyNumber int    `json:"lucky_number"`
	Mood        string `json:"mood"`
}
