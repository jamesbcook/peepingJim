AUTHOR=-X github.com/jamesbcook/peepingJim.Version=4.0.0
VERSION=-X "github.com/jamesbcook/peepingJim.Author=James Cook <@_jbcook>"

OUTPUT_DIR=bin
setup:
	mkdir -p $(OUTPUT_DIR)

linux: setup
	GOOS=linux GOARCH=amd64 go build -ldflags '$(AUTHOR) $(VERSION)' -v -o $(OUTPUT_DIR)/peepingJim_linux_amd64 cmd/main.go

osx: setup
	GOOS=darwin GOARCH=amd64 go build -ldflags '$(AUTHOR) $(VERSION)' -v -o $(OUTPUT_DIR)/peepingJim_darwin_amd64 cmd/main.go

docker: setup
	go build -ldflags '$(AUTHOR) $(VERSION)' -v -o $(OUTPUT_DIR)/peepingJim cmd/main.go

release: linux osx

clean:
	rm -r $(OUTPUT_DIR)
