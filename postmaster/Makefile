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

.PHONY ship-ci:
ship-ci:
	@echo "Shipping..."
	@make build
	@make bundle
	@make codesign
	@make notarize-ci
	@make compress

.PHONY build:
build:
	@echo "Building..."
	@mkdir -p Building
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o bin/pillarBox-amd64-darwin ./cmd/main.go
	@CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o bin/pillarBox-arm64-darwin ./cmd/main.go
	@lipo bin/PillarBox-amd64-darwin bin/PillarBox-arm64-darwin -create -output "bin/pillarbox"

.PHONY bundle:
bundle:
	@mkdir -p "bin/Pillar Box.app/Contents/MacOS"
	@mkdir -p "bin/Pillar Box.app/Contents/Resources"
	@mv "bin/pillarbox" "bin/Pillar Box.app/Contents/MacOS/"
	@cp resources/icon.icns "bin/Pillar Box.app/Contents/Resources/"
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
