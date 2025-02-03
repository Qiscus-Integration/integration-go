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

.PHONY: test test/cover

test:
	go test -race ./...

test/cover:
	@THRESHOLD=25.0; \
	go test -race -coverprofile=./coverage.out ./...; \
	COVERAGE_OUTPUT=$$(go tool cover -func=coverage.out); \
	echo "$$COVERAGE_OUTPUT"; \
	COVERAGE=$$(echo "$$COVERAGE_OUTPUT" | awk '/total:/ {print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$COVERAGE < $$THRESHOLD" | bc -l) -eq 1 ]; then \
		echo "Test coverage ($$COVERAGE%) is below the required threshold ($$THRESHOLD%). Please add more tests!"; \
		exit 1; \
	fi

.PHONY: generate
generate:
	go generate ./...
