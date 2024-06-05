import { resolve } from "path";

const listener = {
	listening: false,
	lastConnected: 0,
	retryIn:0,
};

self.addEventListener('install', () => {
	console.log('install');
});

self.addEventListener('activate', () => {
	console.log('activate');

	start().catch(console.error);
});

async function start() {
	while (true) {
		try {
			await startServer();
		} catch (error) {
			console.error('failed to start server, retrying in 5 seconds', error);
		}

		await new Promise((resolve) => setTimeout(resolve, 5000));
	}
}

async function startServer() {
	return new Promise((resolve, reject) => {
		const ws = new WebSocket('ws://localhost:3500/ws');

		ws.onclose = () => {
			console.log('ws closed');
			resolve(void 0);
		};
		ws.onerror = (err) => {
			console.log('ws error', err);
			reject(err);
		};
		ws.onopen = () => console.log('ws open');
		ws.onmessage = (event) => handleMessage(JSON.parse(event.data));
	});
}

function handleMessage(eventData: any) {
	const { code, payload } = eventData;

	switch (code) {
		case 'mfa_code':
			const { code: mfaCode } = payload.mfa_code;
			
			handleMfaCode(mfaCode);
			break;
		default:
			console.error('Unknown message code', code);
	}
}

async function handleMfaCode(code: string) {
	console.log('handleMfaCode', code);

	const [tab] = await chrome.tabs.query({ active: true, lastFocusedWindow: true });
	
	if (!tab || !tab.id) return;

	try {
		await chrome.tabs.sendMessage(tab.id, {
			code: 'mfa_code',
			payload: {
				mfa_code: {
					code,
				},
			},
		});
	} catch (error) {
		console.error(error);
	}
}
