package codeextractor

import "regexp"

var (
	googleAuthCodeRegex = regexp.MustCompile(`(.?)(G\-[0-9]{6})`)
	standardEightRegex  = regexp.MustCompile(`(.?)([0-9]{8})`)
	standardSevenRegex  = regexp.MustCompile(`(.?)([0-9]{7})`)
	standardSixRegex    = regexp.MustCompile(`(.?)([0-9]{6})`)
	standardFiveRegex   = regexp.MustCompile(`(.?)([0-9]{5})`)
	standardFourRegex   = regexp.MustCompile(`(.?)([0-9]{4})`)
	standardThreeRegex  = regexp.MustCompile(`(.?)([0-9]{3})`)
	standardTwoRegex    = regexp.MustCompile(`(.?)([0-9]{2})`)

	splitEightRegex = regexp.MustCompile(`(.?)([0-9]{4}\-[0-9]{4})`)
	splitSixRegex   = regexp.MustCompile(`(.?)([0-9]{3}\-[0-9]{3})`)
	splitFourRegex  = regexp.MustCompile(`(.?)([0-9]{2}\-[0-9]{2})`)

	spacedEightRegex = regexp.MustCompile(`(.?)([0-9]{4}\-[0-9]{4})`)
	spacedSixRegex   = regexp.MustCompile(`(.?)([0-9]{3}\-[0-9]{3})`)
	spacedFourRegex  = regexp.MustCompile(`(.?)([0-9]{2}\-[0-9]{2})`)

	regexControlOrder = []*regexp.Regexp{
		// G-123456
		googleAuthCodeRegex,

		// 12345678
		// 1234-5678
		standardEightRegex,
		splitEightRegex,

		// 1234567
		standardSevenRegex,

		// 123456
		// 123-456
		standardSixRegex,
		splitSixRegex,

		// 12345
		standardFiveRegex,

		// 1234
		// 12-34
		standardFourRegex,
		splitFourRegex,

		// 123
		// 12
		standardThreeRegex,
		standardTwoRegex,

		// 1234 5678
		// 123 456
		// 12 34
		spacedEightRegex,
		spacedSixRegex,
		spacedFourRegex,
	}
)

// ExtractMFACodeFromMessage extracts an MFA code from a message. If no code is found then
// nil is returned.
func ExtractMFACodeFromMessage(message string) (string, error) {
	for _, regex := range regexControlOrder {
		matches := regex.FindStringSubmatch(message)
		if len(matches) < 2 {
			continue
		}

		return matches[2], nil
	}

	return "", nil
}
