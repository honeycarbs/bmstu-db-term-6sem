package model

type Answer struct {
	Word string
	Mark int
}

type Stats struct {
	Word  string
	Stats float64
}

type Poll struct {
	Answer []Answer
}
