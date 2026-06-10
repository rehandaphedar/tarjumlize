package main

type WordRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type Segment struct {
	Text      string    `json:"t"`
	WordRange WordRange `json:"word_range"`
}

type Verse struct {
	Text     string    `json:"t"`
	Segments []Segment `json:"segments"`
}

type TemplateData struct {
	VerseKey    string
	Arabic      []string
	Translation string
}

type MediaEntry struct {
	Src string `yaml:"src"`
	As  string `yaml:"as"`
}
