# Go project Makefile

PKG=./...          # default: run on all packages
COVERFILE=coverage.out

.PHONY: all test cover bench clean

# Run all tests (verbose)
test:
	go test -v $(PKG)

# Run tests with race detector
race:
	go test -race -v $(PKG)

# Run tests with coverage
cover:
	go test -coverprofile=$(COVERFILE) $(PKG)
	go tool cover -func=$(COVERFILE)
	go tool cover -html=$(COVERFILE)

# Run benchmarks
bench:
	go test -bench=. -benchmem $(PKG)

# Clean up coverage files
clean:
	rm -f $(COVERFILE)
