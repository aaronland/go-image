GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/transform cmd/transform/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/resize cmd/resize/main.go
