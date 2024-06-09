# postmaster

Postmaster is a macOS application that listens for incoming SMS messages, scans them to catch any 2FA/MFA/OTP codes, and sends them to a chrome extension, or to your clipboard if you toggle the option in the menu bar.

## Getting started

```bash
$ go mod download
$ CGO_ENABLED=1 go run cmd/main.go
```
