default: help

help:                             ## Display this help message
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
		awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'

init:                             ## Install development tools
	rm -rf ./bin
	cd tools && go generate -x -tags=tools

ci-init:                ## Initialize CI environment
	# nothing there yet

format:                           ## Format source code
	bin/gofumpt -l -w .
	bin/goimports -local github.com/percona-platform/dbaas-controller -l -w .
	bin/gci write --section Standard --section Default --section "Prefix(github.com/percona-csalguero/yaml_utils)" .

check:                            ## Run checks/linters for the whole project
	bin/check-license
	bin/go-consistent -pedantic -exclude "tests" ./...
	cd tests && ../bin/go-consistent -pedantic ./...
	bin/golangci-lint run

