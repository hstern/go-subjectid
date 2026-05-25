# go-subjectid maintainer Makefile.
#
# Only the gen-grammars / check-gen targets need invoking by hand;
# day-to-day work uses `go test`, `go vet`, `golangci-lint run` directly.

.PHONY: gen-grammars check-gen test lint vet fmt

# gen-grammars re-runs the ABNF -> Go regex generator and writes
# grammars_gen.go. The generator lives in its own module
# (tools/go.mod) so the runtime library's go.mod stays stdlib-only.
gen-grammars:
	cd tools/genabnf && go run . -out ../../grammars_gen.go
	gofmt -w grammars_gen.go

# check-gen re-runs the generator and fails if grammars_gen.go drifts
# from the ABNF inputs. Wired into CI so a maintainer who edits
# abnf/*.abnf but forgets to regenerate gets a red build instead of
# a silent skew between sources and generated code.
check-gen: gen-grammars
	@if ! git diff --quiet -- grammars_gen.go; then \
		echo "ERROR: grammars_gen.go is out of sync with abnf/ sources."; \
		echo "       Run 'make gen-grammars' and commit the result."; \
		git --no-pager diff -- grammars_gen.go; \
		exit 1; \
	fi

# Convenience aggregations of the per-tool commands CI runs.
test:
	go test -race -shuffle=on ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

fmt:
	gofmt -w .
