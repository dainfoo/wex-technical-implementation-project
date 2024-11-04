# ==================================================================================== #
# WARNING: This Makefile requires a UNIX like OS and a modern x86 64 bits processor.   #
# ==================================================================================== #

# ==================================================================================== #
# VARIABLES - You can override some of these environment variables in the command.     #
# ==================================================================================== #

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOAMD64 ?= $(shell if [ "$(GOARCH)" = "amd64" ]; then echo "v3"; else echo ""; fi)

main_path = ./cmd/main.go
binary_build_directory ?= ./tmp/bin
binary_build_name ?= wex
test_coverage_file_path = ./tmp/coverage.out
air_log_directory = ./tmp/log
air_log_name = air.log

# ==================================================================================== #
# HELPERS - These are helper targets that should not be changed.                       #
# ==================================================================================== #

## help: Display this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## confirm: Prompt for confirmation before proceeding
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## no-dirty: Ensure the working directory is clean
.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"

# ==================================================================================== #
# QUALITY CONTROL - These targets are used to ensure the quality of the codebase.      #
# ==================================================================================== #

## test: Run all tests
.PHONY: test
test:
	CGO_ENABLED=1 go test -v -race -buildvcs ./...

## test/cover: Run all tests and display coverage
.PHONY: test/cover
test/cover:
	CGO_ENABLED=1 \
		go test -v -race -buildvcs -coverprofile=${test_coverage_file_path} ./...
	go tool cover -html=${test_coverage_file_path}

## audit: Run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# ==================================================================================== #
# DEVELOPMENT - These targets are used to aid in the development of the application.   #
# ==================================================================================== #

## tidy: Tidy module files and format Go source files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## pre-build: Prepare for building the application
.PHONY: pre-build
pre-build: tidy
	mkdir -p ${binary_build_directory} ${air_log_directory}

## build: Build the application
.PHONY: build
build: pre-build
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GOAMD64=${GOAMD64} \
		go build -o ${binary_build_directory}/${binary_build_name} \
		-trimpath -ldflags="-w -s" ${main_path}

## run/dev: Run the application with live reloading
.PHONY: run/dev
run/dev: pre-build
	go run github.com/air-verse/air@latest \
		--root "." \
		--tmp_dir "tmp" \
		--build.bin "${binary_build_directory}/${binary_build_name}" \
		--build.cmd "make build" \
		--build.delay "1000" \
		--build.include_ext "go" \
		--build.include_file "go.mod" \
		--build.exclude_dir "tmp" \
		--build.exclude_regex "_test.go" \
		--build.exclude_unchanged "false" \
		--build.follow_symlink "false" \
		--build.kill_delay "1s" \
		--build.log "${air_log_directory}/${air_log_name}" \
		--build.poll "false" \
		--build.rerun "false" \
		--build.send_interrupt "true" \
		--build.stop_on_error "false" \
		--log.main_only "false" \
		--log.time "true" \
		--misc.clean_on_exit "false" \
		--screen.clear_on_rebuild "false" \
		--screen.keep_scroll "true"

## clean: Remove temporary files and directories
.PHONY: clean
clean:
	rm -rf ./tmp

# ==================================================================================== #
# OPERATIONS - These targets are intended to be used on operational tasks.             #
# ==================================================================================== #

## pull: Pull changes from the remote Git repository
.PHONY: pull
pull: confirm
	git pull

## push: Push changes to the remote Git repository
.PHONY: push
push: confirm audit no-dirty
	git push
