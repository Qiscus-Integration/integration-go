.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

.PHONY: build
build:
	go build -o /tmp/bin/api cmd/main.go

.PHONY: run
run: build
	/tmp/bin/api $(app)

.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "/tmp/bin/api $(app)" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"