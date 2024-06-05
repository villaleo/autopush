mount:
	@sudo mv bin/autopush /usr/local/bin

deps:
	@/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
	@brew install go

build: deps
	@go build -o bin/autopush main.go

run: build mount
	autopush
