// <input id="twofactorcode" name="twofactorcode" class="required invalid" inputmode="numeric" maxlength="12" value="" data-validation-settings="{&quot;required&quot;:true,&quot;isNumericString&quot;:true}" aria-describedby="twofactorcode_message" aria-required="true" data-invalid="true" aria-invalid="true">
//  <input aria-label="Enter an OTP Code" aria-invalid="false" aria-required="false" autocomplete="off" id="PHONE_SMS_OTP-0" inputmode="numeric" pattern="\d*" placeholder="" type="text" class="ag gc gd ge gf gg gh gi gj gk bu au gl ei gm by gn go gp gq fz bw g0 bd b9 gr gs gt" value="">

const orderedInputSelectors = [
	'input[inputmode="numeric"]',
	'input[type="text"]',
];

chrome.runtime.onMessage.addListener(function(message, sender, sendResponse) {
	const { code, payload } = message;

	if (code !== 'mfa_code')
		return;

	const { code: mfaCode } = payload.mfa_code;

	const input = findInput();

	if (!input) return;

	input.focus();
	input.value = mfaCode;
});

function findInput() {
	for (const selector of orderedInputSelectors) {
		const input = document.querySelector(selector) as HTMLInputElement;

		if (input)
			return input;
	}

	return null;
}
