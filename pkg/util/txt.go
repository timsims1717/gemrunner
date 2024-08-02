package util

import goaway "github.com/TwiN/go-away"

var ProfanityDetector = goaway.NewProfanityDetector().
	WithSanitizeLeetSpeak(false).
	WithSanitizeSpecialCharacters(false).
	WithSanitizeAccents(false)
