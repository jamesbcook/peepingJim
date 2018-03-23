AUTHOR=-X github.com/jamesbcook/peepingJim.Version=3.1.1
VERSION=-X "github.com/jamesbcook/peepingJim.Author=James Cook <@_jbcook>"

OUTPUT_DIR=bin
PHANTOMJS=phantomjs-2.1.1
LINUX=linux-x86_64
OSX=macosx
CAPTUREJS_URL=https://raw.githubusercontent.com/jamesbcook/peepingJim/master/capture.js
DOWNLOAD_URL=https://bitbucket.org/ariya/phantomjs/downloads

setup:
	mkdir -p $(OUTPUT_DIR)

captureJS: 
	wget -O $(OUTPUT_DIR)/capture.js $(CAPTUREJS_URL)

phantomjs-linux:
	wget -O $(OUTPUT_DIR)/$(PHANTOMJS)-$(LINUX).tar.bz2 $(DOWNLOAD_URL)/$(PHANTOMJS)-$(LINUX).tar.bz2
	tar -xvJf $(OUTPUT_DIR)/$(PHANTOMJS)-$(LINUX).tar.bz2 -C $(OUTPUT_DIR)/
	cp $(OUTPUT_DIR)/$(PHANTOMJS)-$(LINUX)/bin/phantomjs $(OUTPUT_DIR)/
	rm $(OUTPUT_DIR)/$(PHANTOMJS)-$(LINUX).tar.bz2
	rm -r $(OUTPUT_DIR)/$(PHANTOMJS)-$(LINUX)

phantomjs-osx:
	wget -O $(OUTPUT_DIR)/$(PHANTOMJS)-$(OSX).zip $(DOWNLOAD_URL)/$(PHANTOMJS)-$(OSX).zip
	unzip $(OUTPUT_DIR)/$(PHANTOMJS)-$(OSX).zip -d $(OUTPUT_DIR)/
	cp $(OUTPUT_DIR)/$(PHANTOMJS)-$(OSX)/bin/phantomjs $(OUTPUT_DIR)/
	rm $(OUTPUT_DIR)/$(PHANTOMJS)-$(OSX).zip
	rm -r $(OUTPUT_DIR)/$(PHANTOMJS)-$(OSX)

linux: setup captureJS phantomjs-linux
	GOOS=linux GOARCH=amd64 go build -ldflags '$(AUTHOR) $(VERSION)' -v -o $(OUTPUT_DIR)/peepingJim_linux_amd64 cmd/main.go

osx: setup captureJS phantomjs-osx
	GOOS=darwin GOARCH=amd64 go build -ldflags '$(AUTHOR) $(VERSION)' -v -o $(OUTPUT_DIR)/peepingJim_darwin_amd64 cmd/main.go

release: linux osx

clean:
	rm -r $(OUTPUT_DIR)
