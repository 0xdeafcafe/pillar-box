build:
	go build -o bin/pillar-box-server cmd/main.go
	# gon ./gon.hcl
	rm -rf ./bin/PillarBox.app
	mkdir ./bin/PillarBox.app
	mkdir ./bin/PillarBox.app/Contents
	mkdir ./bin/PillarBox.app/Contents/MacOS
	mkdir ./bin/PillarBox.app/Contents/Resources
	cp ./resources/Info.plist ./bin/PillarBox.app/Contents/Info.plist
	cp ./resources/icon.icns ./bin/PillarBox.app/Contents/Resources/icon.icns
	cp ./bin/pillar-box-server ./bin/PillarBox.app/Contents/MacOS/PillarBox.app
