.PHONY clean:
clean:
	@echo "Cleaning..."
	@rm -rf bin

.PHONY ship:
ship:
	@make clean
	@echo "Shipping..."
	@make build
	@make bundle
	@make codesign
	@make notarize

.PHONY ship-beta:
ship-beta:
	@make clean
	@echo "Shipping..."
	@make build
	@make bundle-beta
	@make codesign
	@make notarize

.PHONY ship-ci:
ship-ci:
	@echo "Shipping..."
	@make build-ci
	@make bundle
	@make codesign
	@make notarize-ci
	@make compress

.PHONY ship-ci-beta:
ship-ci-beta:
	@echo "Shipping..."
	@make build-ci
	@make bundle-beta
	@make codesign
	@make notarize-ci
	@make compress

.PHONY build:
build:
	@echo "Building..."
	@mkdir -p bin
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o bin/pillarbox-amd64-darwin ./cmd/main.go
	@CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o bin/pillarbox-arm64-darwin ./cmd/main.go
	@lipo bin/pillarbox-amd64-darwin bin/pillarbox-arm64-darwin -create -output "bin/pillarbox"

.PHONY build-ci:
build-ci:
	@echo "Building..."
	@mkdir -p bin
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -v -ldflags="-X github.com/0xdeafcafe/pillar-box/server/internal/updater/updater.Version=$PB_VERSION" -o bin/pillarbox-amd64-darwin ./cmd/main.go
	@CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -v -ldflags="-X github.com/0xdeafcafe/pillar-box/server/internal/updater/updater.Version=$PB_VERSION" -o bin/pillarbox-arm64-darwin ./cmd/main.go
	@lipo bin/pillarbox-amd64-darwin bin/pillarbox-arm64-darwin -create -output "bin/pillarbox"

.PHONY bundle:
bundle:
	@mkdir -p "bin/Pillar Box.app/Contents/MacOS"
	@mkdir -p "bin/Pillar Box.app/Contents/Resources"
	@mv "bin/pillarbox" "bin/Pillar Box.app/Contents/MacOS/"
	@cp resources/icon.icns "bin/Pillar Box.app/Contents/Resources/"
	@cp resources/Info.plist "bin/Pillar Box.app/Contents/"

.PHONY bundle-beta:
bundle-beta:
	@mkdir -p "bin/Pillar Box.app/Contents/MacOS"
	@mkdir -p "bin/Pillar Box.app/Contents/Resources"
	@mv "bin/pillarbox" "bin/Pillar Box.app/Contents/MacOS/"
	@cp resources/icon-beta.icns "bin/Pillar Box.app/Contents/Resources/icon.icns"
	@cp resources/Info.plist "bin/Pillar Box.app/Contents/"

.PHONY codesign:
codesign:
	@echo "Signing..."
	@codesign --deep --force --verify --verbose \
		--sign "$(APPLE_CODESIGN_IDENTITY)" \
		--options runtime \
		--timestamp \
		"bin/Pillar Box.app"
	@codesign --verify --verbose=4 "bin/Pillar Box.app"

.PHONY notarize-ci:
notarize-ci:
	@echo "Notarizing..."
	@zip -r "bin/Pillar Box.zip" "bin/Pillar Box.app"
	@xcrun notarytool store-credentials "pillar-box" \
		--apple-id "$(APPLE_ID)" \
		--team-id "$(APPLE_TEAM_ID)" \
		--password "$(APPLE_APP_SPECIFIC_PASSWORD)"
	@xcrun notarytool submit "bin/Pillar Box.zip" \
		--keychain-profile "pillar-box" \
		--wait
	@xcrun stapler staple "bin/Pillar Box.app"
	@spctl --assess --verbose "bin/Pillar Box.app"

.PHONY notarize:
notarize:
	@echo "Notarizing..."
	@zip -r "bin/Pillar Box.zip" "bin/Pillar Box.app"
	@xcrun notarytool submit "bin/Pillar Box.zip" \
		--keychain-profile "alex.forbes-reed" \
		--wait
	@xcrun stapler staple "bin/Pillar Box.app"
	@spctl --assess --verbose "bin/Pillar Box.app"

.PHONY compress:
compress:
	@echo "Compressing..."
	@zip -r "bin/Pillar Box.zip" "bin/Pillar Box.app"
