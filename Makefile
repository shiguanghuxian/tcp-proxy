default:
	@echo 'Usage of make: [ build | linux_build | windows_build | clean ]'

build: 
	@go build -o ./bin/tcp-proxy ./

linux_build: 
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/tcp-proxy ./

windows_build: 
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/tcp-proxy.exe ./

clean: 
	@rm -f ./bin/tcp-proxy*

.PHONY: default build linux_build windows_build clean