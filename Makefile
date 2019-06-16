all: init build

init:
	@echo "Build all commands for linux, darwin platform"

build:
	@echo "Build for Linux Plaform(386)"
	@echo "Build cmd/gdrivesql/gdrivesql.go file"
	env GOOS=linux GOARCH=386 go build -o gdrivesql-linux-386 cmd/gdrivesql/gdrivesql.go
	@echo "Build cmd/gdriveauth/gdriveauth.go file"
	env GOOS=linux GOARCH=386 go build -o gdriveauth-linux-386 cmd/gdriveauth/gdriveauth.go
	@echo "Build cmd/cleanup/cleanup.go file"
	env GOOS=linux GOARCH=386 go build -o cleanup-linux-386 cmd/cleanup/cleanup.go

	@echo "Build for Darwin Plaform(386)"
	@echo "Build cmd/gdrivesql/gdrivesql.go file"
	env GOOS=darwin GOARCH=386 go build -o gdrivesql-386 cmd/gdrivesql/gdrivesql.go
	@echo "Build cmd/gdriveauth/gdriveauth.go file"
	env GOOS=darwin GOARCH=386 go build -o gdriveauth-386 cmd/gdriveauth/gdriveauth.go
	@echo "Build cmd/cleanup/cleanup.go file"
	env GOOS=darwin GOARCH=386 go build -o cleanup-386 cmd/cleanup/cleanup.go