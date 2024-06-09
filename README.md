# Pillar Box

Pillar Box recreates the "AutoFill for SMS codes" feature that is available in Safari on
macOS. It consists of a macOS application which will monitor your messages for incoming
SMS codes, and a Chrome extension which will communicate with the macOS application to
listen for SMS codes and automatically fill them in.

## Key features

- Automatically fill in SMS codes from your messages
- No need to copy and paste codes from messages, or switch between apps
- Secure, nothing leaves your machine

## Running locally

### macOS application

```bash
$ cd postmaster
$ go run cmd/main.go
```

### Chromium extension

1. Open `chrome://extensions/` in your browser (tested with Chrome and Arc)
2. Enable "Developer mode" in the top right corner
3. Click "Load unpacked" and select the `extension` directory

