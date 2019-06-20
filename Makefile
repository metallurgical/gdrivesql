all: init build

init:
	@echo "Build all commands for linux, darwin platform"

build:
	@echo "Build for Linux Plaform(386)"
	@echo "Build cmd/gdrivesql/gdrivesql.go file"
	env GOOS=linux GOARCH=386 go build -o gdrivesql-linux-386 cmd/gdrivesql/gdrivesql.go
	@echo "Build cmd/gdriveauth/gdriveauth.go file"
	env GOOS=linux GOARCH=386 go build -o gdriveauth-linux-386 cmd/gdriveauth/gdriveauth.go
	@echo "Build cmd/gdriveclean/gdriveclean.go file"
	env GOOS=linux GOARCH=386 go build -o gdriveclean-linux-386 cmd/gdriveclean/gdriveclean.go

	@echo "Build for Linux Plaform(amd64)"
	@echo "Build cmd/gdrivesql/gdrivesql.go file"
	env GOOS=linux GOARCH=amd64 go build -o gdrivesql-linux-amd64 cmd/gdrivesql/gdrivesql.go
	@echo "Build cmd/gdriveauth/gdriveauth.go file"
	env GOOS=linux GOARCH=amd64 go build -o gdriveauth-linux-amd64 cmd/gdriveauth/gdriveauth.go
	@echo "Build cmd/gdriveclean/gdriveclean.go file"
	env GOOS=linux GOARCH=amd64 go build -o gdriveclean-linux-amd64 cmd/gdriveclean/gdriveclean.go

	@echo "Build for Darwin Plaform(386)"
	@echo "Build cmd/gdrivesql/gdrivesql.go file"
	env GOOS=darwin GOARCH=386 go build -o gdrivesql-386 cmd/gdrivesql/gdrivesql.go
	@echo "Build cmd/gdriveauth/gdriveauth.go file"
	env GOOS=darwin GOARCH=386 go build -o gdriveauth-386 cmd/gdriveauth/gdriveauth.go
	@echo "Build cmd/gdriveclean/gdriveclean.go file"
	env GOOS=darwin GOARCH=386 go build -o gdriveclean-386 cmd/gdriveclean/gdriveclean.go

	@echo "Build for Darwin Plaform(amd64)"
	@echo "Build cmd/gdrivesql/gdrivesql.go file"
	env GOOS=darwin GOARCH=amd64 go build -o gdrivesql-amd64 cmd/gdrivesql/gdrivesql.go
	@echo "Build cmd/gdriveauth/gdriveauth.go file"
	env GOOS=darwin GOARCH=amd64 go build -o gdriveauth-amd64 cmd/gdriveauth/gdriveauth.go
	@echo "Build cmd/gdriveclean/gdriveclean.go file"
	env GOOS=darwin GOARCH=amd64 go build -o gdriveclean-amd64 cmd/gdriveclean/gdriveclean.go