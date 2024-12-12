package codeextractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractMFACodeFromMessage_TestCases(t *testing.T) {
	type TestCase struct {
		Message string
		Code    string
	}

	testCases := []TestCase{{
		Message: `Uw sms-code is: 481243. Deze vervalt over 20 minuten.`,
		Code:    "481243",
	}, {
		Message: `Your Uber code is 7791. Never share this code.`,
		Code:    "7791",
	}, {
		Message: `Uw sms-code is: 296891`,
		Code:    "296891",
	}, {
		Message: `263810 is your verification code for Outside: Shared Calendar App.`,
		Code:    "263810",
	}, {
		Message: `Uw SMS code is 462838`,
		Code:    "462838",
	}, {
		Message: `NEVER share this One-Time Code: 123456. Amex will never call to ask for it. If released to someone or not requested, call us using Contact Us on Amex website`,
		Code:    "123456",
	}, {
		Message: `Amex SafeKey code is 123456 for €13.37 transaction attempt at Bol.com for Card ending in 69420. Never share this code.`,
		Code:    "123456",
	}, {
		Message: `Your OpenTable verification code is: 207734. This code will expire in 10 minutes. Don't share this code with anyone; our employees will never ask for the code.`,
		Code:    "207734",
	}, {
		Message: `354352 is your verification code for Telegram Messenger.`,
		Code:    "354352",
	}, {
		Message: `Glovo code: 3851. Valid for 3 minutes.`,
		Code:    "3851",
	}, {
		Message: `Your Raya verification code is: 26576`,
		Code:    "26576",
	}, {
		Message: `Your app verification code is: 717949. Don't share this code with anyone; our employees will never ask for the code.`,
		Code:    "717949",
	}, {
		Message: `To finish registering Authy click: authy://register/12-345-678-9090/123456 or manually enter: 123456`,
		Code:    "123456",
	}, {
		Message: `Your locker: 0962`,
		Code:    "0962",
	}, {
		Message: `Your access code: 726259\n\nOpen your locker:\nwww.yourlocker.nl/en/open/tantetoos/0962/726259`,
		Code:    "726259",
	}, {
		Message: `Snapchat Login Code: 351304. Snapchat Support will not ask for this code. Do not share it with anyone.`,
		Code:    "351304",
	}, {
		Message: `<#>Your Deliveroo verification code is: 795800\n/tPjtJT5f8o`,
		Code:    "795800",
	}, {
		Message: `Use 123456 as your verification code to confirm your log in to Barclaycard at 12:12:12 on 11/11/2011. Don't recognise this? Call us immediately.`,
		Code:    "123456",
	}, {
		Message: `Uw Infomedics inlog code is: 399516.`,
		Code:    "399516",
	}, {
		Message: `VIVA WALLET CODE: 237892 TO AUTHORISE ACCESS TO YOUR ACCOUNT.\nIP: 104.28.42.75\n\nNOT YOU? CONTACT US IN THE APP.`,
		Code:    "237892",
	}, {
		Message: `6213 Gebruik deze code om in te loggen in jouw Getir account. Deel deze code met niemand en gebruik hem alleen om je account te activeren.`,
		Code:    "6213",
	}, {
		Message: `Companies House Notifications: 777353 is your security code. Your code expires in 30 minutes.`,
		Code:    "777353",
	}, {
		Message: `Uw code is: 4616. Bedankt.`,
		Code:    "4616",
	}, {
		Message: `Your Coinbase verification code is: 2287646. Don't share this code with anyone; our employees will never ask for the code.`,
		Code:    "2287646",
	}, {
		Message: `G-743499 is your Google verification code.`,
		Code:    "G-743499",
	}, {
		Message: `Your verification code is 903504`,
		Code:    "903504",
	}, {
		Message: `Your Uber code is 2302. Never share this code. Reply STOP ALL to +44 7903 561836 to unsubscribe.`,
		Code:    "2302",
	}}

	for _, testCase := range testCases {
		code, err := ExtractMFACodeFromMessage(testCase.Message)

		assert.NoError(t, err)
		assert.NotNil(t, code)
		assert.Equal(t, testCase.Code, *code)
	}
}
