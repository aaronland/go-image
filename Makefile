GOMOD=vendor

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/transform cmd/transform/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/resize cmd/resize/main.go
