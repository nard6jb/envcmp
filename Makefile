.PHONY: build test lint clean

BINARY := envcmp
CMD     := ./cmd/envcmp

build:
	go build -o $(BINARY) $(CMD)

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)

run-diff: build
	./$(BINARY) diff $(FILE1) $(FILE2)

run-validate: build
	./$(BINARY) validate $(REF) $(TARGET)
