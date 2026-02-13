package moderation

import "strings"

// profanityList contains common English profanity words used for automated
// content screening. The list intentionally contains offensive language so that
// user-submitted content can be checked against it.
var profanityList = []string{
	"ass",
	"asshole",
	"bastard",
	"bitch",
	"bitches",
	"bollocks",
	"bullshit",
	"cock",
	"cocksucker",
	"crap",
	"cunt",
	"damn",
	"dick",
	"dickhead",
	"douche",
	"douchebag",
	"dumbass",
	"fag",
	"faggot",
	"fuck",
	"fucker",
	"fucking",
	"goddamn",
	"hell",
	"horseshit",
	"jackass",
	"jerkoff",
	"motherfucker",
	"nigga",
	"nigger",
	"piss",
	"prick",
	"pussy",
	"retard",
	"retarded",
	"shit",
	"shithead",
	"slut",
	"son of a bitch",
	"spic",
	"twat",
	"wanker",
	"whore",
	"wtf",
	"stfu",
	"milf",
	"dildo",
	"blowjob",
	"handjob",
	"porn",
}

// ContainsProfanity returns true if the given text contains any word from the
// built-in profanity list. The check is case-insensitive and matches whole
// words by splitting on whitespace and common punctuation boundaries.
func ContainsProfanity(text string) bool {
	lower := strings.ToLower(text)
	for _, word := range profanityList {
		// Check for the word surrounded by word boundaries (space, start/end, punctuation)
		if containsWord(lower, word) {
			return true
		}
	}
	return false
}

// containsWord checks if text contains word as a whole word (not as a substring
// of a larger word). For multi-word phrases (e.g. "son of a bitch") it uses
// simple substring matching.
func containsWord(text, word string) bool {
	if strings.Contains(word, " ") {
		// Multi-word phrase: substring match is fine
		return strings.Contains(text, word)
	}

	// Split on common boundaries to extract individual tokens
	tokens := strings.FieldsFunc(text, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'))
	})
	for _, t := range tokens {
		if t == word {
			return true
		}
	}
	return false
}

// ScanContent examines a title and description for policy violations and
// returns a list of reason strings (e.g. "profanity_in_title"). An empty
// slice means no issues were found.
func ScanContent(title, description string) []string {
	var reasons []string
	if ContainsProfanity(title) {
		reasons = append(reasons, "profanity_in_title")
	}
	if ContainsProfanity(description) {
		reasons = append(reasons, "profanity_in_description")
	}
	return reasons
}
