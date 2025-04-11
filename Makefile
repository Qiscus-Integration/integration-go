.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

.PHONY: build
build:
	go build -o /tmp/bin/app main.go

.PHONY: run
run: build
	/tmp/bin/app $(bin)

.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "/tmp/bin/app $(bin)" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"

.PHONY: test test/coverage

test:
	go test $$(go list ./... | grep -v 'test\|mocks') -race -coverprofile=./coverage.out

test/coverage: test
	@THRESHOLD=25.0; \
	COVERAGE_OUTPUT=$$(go tool cover -func=coverage.out); \
	COVERAGE=$$(echo "$$COVERAGE_OUTPUT" | awk '/total:/ {print $$3}' | sed 's/%//'); \
	echo "Total test coverage: $$COVERAGE%"; \
	if [ $$(awk "BEGIN {if ($$COVERAGE < $$THRESHOLD) print 1; else print 0}") -eq 1 ]; then \
		echo "Test coverage is below the required threshold ($$THRESHOLD%). Please add more tests!"; \
		exit 1; \
	fi

.PHONY: generate
generate:
	go generate ./...
