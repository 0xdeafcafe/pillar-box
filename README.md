# Pillar Box

Pillar Box recreates the "AutoFill for SMS codes" feature that is available in Safari on
macOS. It consists of a macOS application which will monitor your messages for incoming
SMS codes, and a Chrome extension which will communicate with the macOS application to
listen for SMS codes and automatically fill them in.

## Key features

- Automatically fill in SMS codes from your messages
- No need to copy and paste codes from messages, or switch between apps
- Secure, nothing leaves your machine

## Pre-cooked

The easiest way to get started is to download the latest release from the [releases page](releases). Download both files:

- `Pillar Box.zip` - The macOS application
- `pillar-box-chromium-extension.zip` - The Chromium extension

Extract the app from the archive and run it, you may need to go to System Preferences > Privacy & Security > Full Disk Access and allow the app for it to run.

Extract the extension from the archive and navigate to `chrome://extensions/` in your browser. Enable "Developer mode" in the top right corner, click "Load unpacked" and select the extension you just extracted.

That should be it, go to a website that requires an SMS code and you should see the code automatically filled in.

## Running from source

### macOS application

```bash
$ cd postmaster
$ go run cmd/main.go
```

### Chromium extension

1. Open `chrome://extensions/` in your browser (tested with Chrome and Arc)
2. Enable "Developer mode" in the top right corner
3. Click "Load unpacked" and select the `extension` directory

