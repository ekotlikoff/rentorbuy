.PHONY: all

all:
	go run cmd/rentorbuy.go -f data/data.json

interactive:
	go run cmd/rentorbuy.go -i
