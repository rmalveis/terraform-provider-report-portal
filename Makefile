CURRENT_VERSION=0.0.1
TARGET=darwin_amd64
TARGET_LINUX=linux_amd64

default: build

build: fmtcheck build-ci

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

build-ci:
	go build -o bin/terraform-provider-reportportal
fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./provider $(filter-out ./awsproviderlint/go% ./awsproviderlint/README.md ./awsproviderlint/vendor, $(wildcard ./awsproviderlint/*))

test: fmtcheck
	go test ./provider -timeout=5m -parallel=4

deploylocal: build
	@echo "==> Copying the binary into the plugin folder"
	mkdir -p ~/.terraform.d/plugins/raiadrogasil.com/ecommerce/reportportal/${CURRENT_VERSION}/${TARGET}
	cp -f bin/terraform-provider-reportportal ~/.terraform.d/plugins/raiadrogasil.com/ecommerce/reportportal/${CURRENT_VERSION}/${TARGET}/

deploylinux: fmt build
	@echo "==> Copying the binary into the plugin folder"
	mkdir -p ~/.terraform.d/plugins/acme.com/reportportal/${CURRENT_VERSION}/${TARGET_LINUX}
	cp -f bin/terraform-provider-reportportal ~/.terraform.d/plugins/acme.com/reportportal/${CURRENT_VERSION}/${TARGET_LINUX}/

