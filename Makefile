undefine GOFLAGS

GOLANGCI_LINT_VERSION?=v2.5.0
GOTESTSUM_VERSION?=v1.13.0
GO_TEST?=go run gotest.tools/gotestsum@$(GOTESTSUM_VERSION) --format testname --
TIMEOUT := "60m"

ifeq ($(shell command -v podman 2> /dev/null),)
	RUNNER=docker
else
	RUNNER=podman
endif

# if the golangci-lint steps fails with one of the following error messages:
#
#   directory prefix . does not contain main module or its selected dependencies
#
#   failed to initialize build cache at /root/.cache/golangci-lint: mkdir /root/.cache/golangci-lint: permission denied
#
# you probably have to fix the SELinux security context for root directory plus your cache
#
#   chcon -Rt svirt_sandbox_file_t .
#   chcon -Rt svirt_sandbox_file_t ~/.cache/golangci-lint
lint:
	mkdir -p ~/.cache/golangci-lint/$(GOLANGCI_LINT_VERSION)
	$(RUNNER) run -t --rm \
		-v $(shell pwd):/app \
		-v ~/.cache/golangci-lint/$(GOLANGCI_LINT_VERSION):/root/.cache \
		-w /app \
		golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run -v --max-same-issues 50
.PHONY: lint

format:
	gofmt -w -s $(shell pwd)
.PHONY: format

unit:
	$(GO_TEST) -shuffle on ./...
.PHONY: unit

coverage:
	$(GO_TEST) -shuffle on -covermode count -coverprofile cover.out -coverpkg=./... ./...
.PHONY: coverage
