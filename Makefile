# go-subjectid maintainer Makefile.
#
# Convenience wrappers around the per-tool commands CI runs;
# day-to-day work can also call go test / go vet / golangci-lint
# directly.

.PHONY: test lint vet fmt tidy ci

test:
	go test -race -shuffle=on ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy

# ci runs everything CI runs, locally, in order.
ci: vet tidy test lint
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "gofmt drift:"; gofmt -l .; exit 1; \
	fi
