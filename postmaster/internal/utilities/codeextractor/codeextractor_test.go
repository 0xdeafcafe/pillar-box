package codeextractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractCodesEachMessage(t *testing.T) {
	testMessages := []struct {
		name    string
		message string
		want    []string
		wantErr bool
	}{
		{
			name:    "Msg1 - Basic numeric code",
			message: "Uw sms-code is: 205095. Deze vervalt over 20 minuten.",
			want:    []string{"205095"},
		},
		{
			name:    "Msg2 - One-Time Code (Amex)",
			message: `NEVER share this One-Time Code: 650630. Amex will never call to ask for it. ...`,
			want:    []string{"650630"},
		},
		{
			name:    "Msg3 - Another numeric code",
			message: "Uw sms-code is: 592411. Deze vervalt over 20 minuten.",
			want:    []string{"592411"},
		},
		{
			name:    "Msg4 - Code with dash (524-504)",
			message: "Your DigiD SMS code to log into Mijn OHRA Zorgverzekering is: 524-504",
			want:    []string{"524504"},
		},
		{
			name:    "Msg5 - 6-digit code (013034)",
			message: "Your Suspicious Antwerp verification code is: 013034",
			want:    []string{"013034"},
		},
		{
			name:    "Msg6 - Deliveroo code 979700",
			message: "<#>Your Deliveroo verification code is: 979700\n/tPjtJT5f8o",
			want:    []string{"979700"},
		},
		{
			name:    "Msg7 - Uber code 1808",
			message: "Your Uber code is 1808. Never share this code.",
			want:    []string{"1808"},
		},
		{
			name:    "Msg8 - Amex again (650630 repeated)",
			message: `NEVER share this One-Time Code: 650630. ...`,
			want:    []string{"650630"},
		},
		{
			name:    "Msg9 - Jumbo code 972905",
			message: "972905 is je sms-code om verder te gaan bij Jumbo.",
			want:    []string{"972905"},
		},
		{
			name:    "Msg10 - Another Uber code 7866",
			message: "Your Uber code is 7866. Never share this code.",
			want:    []string{"7866"},
		},
		{
			name:    "Msg11 - Shop verification code 350703",
			message: "350703 is your Shop verification code",
			want:    []string{"350703"},
		},
		{
			name:    "Msg12 - Amex SafeKey KLM 346020",
			message: "Amex SafeKey verificatiecode is 346020 voor €2.916,24 bij KLM ...",
			want:    []string{"346020"},
		},
		{
			name:    "Msg13 - Amex SafeKey Apple 932857",
			message: "Amex SafeKey verificatiecode is 932857 voor €568,00 bij Apple ...",
			want:    []string{"932857"},
		},
		{
			name:    "Msg14 - Uw SMS code 380000",
			message: "Uw SMS code is 380000",
			want:    []string{"380000"},
		},
		{
			name:    "Msg15 - OpenTable code 226044",
			message: "Uw OpenTable-verificatiecode is: 226044.. Deze code verloopt over 10 minuten...",
			want:    []string{"226044"},
		},
		{
			name:    "Msg16 - Amex SafeKey 621740",
			message: "Amex SafeKey code is 621740 for €411.34 transaction attempt at KLM ...",
			want:    []string{"621740"},
		},
		{
			name:    "Msg17 - Apple Pay code 402685",
			message: "Your one-time verification code to add your Amex Card to Apple Pay is 402685...",
			want:    []string{"402685"},
		},
		{
			name:    "Msg18 - Tikkie code 8890",
			message: "Tikkie code: 8890\nDeze code is 5 minuten geldig.",
			want:    []string{"8890"},
		},
		{
			name:    "Msg19 - SafeKey One Time Code 437951",
			message: "437951 is your all numeric SafeKey One Time Code...",
			want:    []string{"437951"},
		},
		{
			name:    "Msg20 - DICE verification code 5845",
			message: "Your DICE verification code is: 5845",
			want:    []string{"5845"},
		},
		{
			name:    "Msg21 - Gett code 204065",
			message: "Gett account confirmation code: 204065",
			want:    []string{"204065"},
		},
		{
			name:    "Msg22 - Coinbase code 1203227",
			message: "Your Coinbase verification code is: 1203227. Don't share...",
			want:    []string{"1203227"},
		},
		{
			name:    "Msg23 - Mixpanel code 2394644",
			message: "Your Mixpanel code is 2394644",
			want:    []string{"2394644"},
		},
		{
			name:    "Msg24 - Mailchimp code 964933",
			message: "Your Mailchimp Two Factor Auth verification code is: 964933",
			want:    []string{"964933"},
		},
		{
			name:    "Msg25 - Tesco code 838123",
			message: "838123 is your Tesco authentication code.\n@tesco.com #838123",
			want:    []string{"838123"},
		},
		{
			name:    "Msg26 - Stripe code 214-576",
			message: "Your Stripe verification code is: 214-576. Don't share this code...",
			want:    []string{"214576"},
		},
		{
			name:    "Msg27 - Google code G-089350",
			message: "G-089350 is your Google verification code.",
			want:    []string{"089350"},
		},
		{
			name:    "Msg28 - Stripe code 473-293",
			message: "Your Stripe verification code is: 473-293...",
			want:    []string{"473293"},
		},
		{
			name:    "Msg29 - Another Stripe code 913-170",
			message: "Your verification code for Stripe is 913-170...",
			want:    []string{"913170"},
		},
		{
			name:    "Msg30 - Twitter login code 940326",
			message: "940326 is your Twitter login code. Don't reply...",
			want:    []string{"940326"},
		},
		// You can add more test entries if needed...
	}

	for _, tm := range testMessages {
		t.Run(tm.name, func(t *testing.T) {
			got, err := ExtractCodes(tm.message)

			if tm.wantErr {
				// We expect an error in this case
				assert.Error(t, err, "Expected an error but got none")
				return
			} else {
				// We do not expect an error
				assert.NoError(t, err, "Did not expect an error but got one")
			}

			// Verify that each expected code is contained in the result
			for _, w := range tm.want {
				assert.Contains(t, got, w, "Expected code %q not found in result: %v", w, got)
			}
		})
	}
}
