package main

import (
	"log"

	qul "git.sr.ht/~rehandaphedar/genanki-go-utils/v2/pkg/qul"
)

func renderSegment(wordIndex qul.WordIndex, verseKey string, from, to int) []string {
	return renderRange(wordIndex, qul.Source{Key: verseKey, From: from, To: to})
}

func renderRange(wordIndex qul.WordIndex, source qul.Source) []string {
	words := wordIndex.VerseWords[source.Key]

	if source.To >= len(words) {
		log.Printf("invalid range %+v, silently fixing", source)
		source.To = len(words)
	}

	if (source.To + 1) == len(words) {
		source.To++
	}

	var parts []string
	for i := source.From - 1; i < source.To; i++ {
		parts = append(parts, words[i])
	}
	return parts
}
