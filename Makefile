# go-subjectid maintainer Makefile.
#
# Convenience wrappers around the per-tool commands CI runs;
# day-to-day work can also call go test / go vet / golangci-lint
# directly.

ABNF_SRCS := $(shell find . -name '*.abnf' -not -path './.*')
ABNF_REX  := $(ABNF_SRCS:.abnf=.rex)

.PHONY: test lint vet fmt tidy ci generate

generate: $(ABNF_REX)

# Regenerate %.rex from %.abnf via pap (github.com/pandatix/go-abnf).
# pap quirks worked around in this rule:
#   - rejects standalone ";" comment lines and blank lines (strip both)
#   - requires CRLF line endings (RFC 5234 strict) (sed appends \r)
#   - appends a trailing \n to stdout, which //go:embed keeps and Go's
#     regexp interprets as a required literal newline (tr -d strips it)
#
# Convention: the top-level rule name is "<basename>-url" — e.g. did-url.abnf
# is rooted at "did-url". Adjust below per-file if a grammar ever ships
# with a different root rule.
%.rex: %.abnf
	grep -v '^;' $< | grep -v '^$$' | sed 's/$$/\r/' \
		| pap regex --rulename $(notdir $*) \
		| tr -d '\n\r' > $@

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

clean:
	rm -f $(ABNF_REX)