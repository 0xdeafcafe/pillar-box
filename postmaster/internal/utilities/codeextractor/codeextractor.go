package codeextractor

import (
	"errors"
	"regexp"
	"sort"
	"strings"
)

const (
	backwardsContextWindow = 60
)

var (
	codeRegexPattern = regexp.MustCompile(`(?i)(G-\d{6}|\d{3,4}-\d{3,4}|\b\d{4,7}\b)`)

	// contextIndicators is a set of phrases that, if found near the code, increase the
	// "likelihood" score for that code.
	contextIndicators = []string{
		"one-time code",
		"one time code",
		"verification code",
		"sms-code",
		"sms code",
		"code is",
		"tfa code",
		"2fa code",
		"tikkie code",
		"safekey code",
		"deliveroo verification code",
		"uber code",
		"dice verification code",
		"gett account confirmation code",
		"mixpanel code",
		"coinbase verification code",
		"mailchimp two factor auth verification code",
		"stripe verification code",
		"google verification code",
		"twitter login code",
		"tesco authentication code",
	}

	ErrNoCodesFound = errors.New("no codes found")
)

// codeHit holds intermediate data about one "found code"
type codeHit struct {
	// code contains the discovered code itself
	code string

	// index is the position in the text where the code was found in the text
	index int

	// score is a computed "likelihood" score for this code
	score int
}

// ExtractCodes attempts to find all 2FA codes in the provided text,
// ranks them by "likelihood", removes duplicates, and returns them
// in descending order of likelihood.
func ExtractCodes(text string) ([]string, error) {
	text = strings.TrimSpace(text)
	matchIndexes := codeRegexPattern.FindAllStringIndex(text, -1)
	if len(matchIndexes) == 0 {
		return nil, errors.New("no codes found")
	}

	var codeHits []codeHit
	for _, m := range matchIndexes {
		raw := text[m[0]:m[1]]

		code := sanitizeCode(raw)
		if code == "" {
			continue
		}

		// Build the codeHit
		ch := codeHit{
			code:  code,
			index: m[0],
			score: 0, // will compute next
		}

		codeHits = append(codeHits, ch)
	}

	if len(codeHits) == 0 {
		return nil, ErrNoCodesFound
	}

	// Score each codeHit based on textual context around it.
	// We'll look ~60 characters before the code’s position for any context indicators.
	for i := range codeHits {
		codeHits[i].score = computeContextScore(text, codeHits[i].index, contextIndicators)
	}

	// Sort by score DESC, then by index ASC (if you want earlier-located codes to break ties).
	sort.SliceStable(codeHits, func(i, j int) bool {
		if codeHits[i].score == codeHits[j].score {
			return codeHits[i].index < codeHits[j].index
		}
		return codeHits[i].score > codeHits[j].score
	})

	// Remove duplicates while preserving order
	uniqueOrdered := make([]string, 0, len(codeHits))
	seen := make(map[string]bool)
	for _, ch := range codeHits {
		if !seen[ch.code] {
			seen[ch.code] = true
			uniqueOrdered = append(uniqueOrdered, ch.code)
		}
	}

	// If after deduping, none remain, error
	if len(uniqueOrdered) == 0 {
		return nil, ErrNoCodesFound
	}

	return uniqueOrdered, nil
}

// extractCodeByPattern extracts the raw code (digits only) from the match depending on
// which pattern was used.
func extractCodeByPattern(text string, matchIdx []int, pattern string) (string, int) {
	// We'll re-run the same regex with submatches so we can check them more easily.
	re := regexp.MustCompile(pattern)
	fullMatch := text[matchIdx[0]:matchIdx[1]] // The raw matched substring
	startPos := matchIdx[0]

	switch {
	// G-xxxxxx pattern
	case pattern == `G-(\d{6})`:
		// e.g. "G-089350"
		// group 1 is the digits. So submatch #1 is what we want as digits.
		sub := re.FindStringSubmatch(fullMatch)
		if len(sub) == 2 {
			return sub[1], startPos // remove "G-"
		}

	// e.g. "524-504" or "913-170"
	case pattern == `(\d{3,4})-(\d{3,4})`:
		sub := re.FindStringSubmatch(fullMatch)
		// sub[1] = first group of digits
		// sub[2] = second group
		if len(sub) == 3 {
			return sub[1] + sub[2], startPos
		}

	// plain digits (4 to 7 digits)
	case pattern == `\b(\d{4,7})\b`:
		sub := re.FindStringSubmatch(fullMatch)
		if len(sub) == 2 {
			return sub[1], startPos
		}
	}

	return "", 0
}

// computeContextScore checks for known "trigger phrases" near the code’s location
// and assigns points if found. You can tweak the distance or logic as needed.
func computeContextScore(text string, codePos int, indicators []string) int {
	score := 0

	// figure out the start index for context scanning
	start := codePos - backwardsContextWindow
	if start < 0 {
		start = 0
	}

	vicinity := strings.ToLower(text[start:codePos])

	for _, ind := range indicators {
		if strings.Contains(vicinity, strings.ToLower(ind)) {
			score += 5
		}
	}
	return score
}

// sanitizeCode removes "G-" prefix if present, and also removes any dashes, leaving only
// digits.
func sanitizeCode(raw string) string {
	raw = strings.ToUpper(raw)

	// Remove "G-" prefix if found
	if strings.HasPrefix(raw, "G-") {
		raw = raw[2:] // skip "G-"
	}

	// remove dashes
	raw = strings.ReplaceAll(raw, "-", "")

	// raw should now be digits only, optionally we can check if it's all digits:
	reDigits := regexp.MustCompile(`^\d+$`)
	if reDigits.MatchString(raw) {
		return raw
	}

	return ""
}
