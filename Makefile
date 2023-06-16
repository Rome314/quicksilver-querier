build:
	go mod tidy
	go build -o quickdump ./cmd
