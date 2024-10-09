PACKAGES=$(shell go list ./... | grep -v '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin

DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf

export GO111MODULE = on

DEB_BIN_DIR ?= /usr/local/bin
DEB_LIB_DIR ?= /usr/lib

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
	build_tags += ledger
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
empty = $(whitespace) $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(empty),$(comma),$(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=odinchain \
	-X github.com/cosmos/cosmos-sdk/version.AppName=odind \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags_comma_sep)" -ldflags '$(ldflags)'

all: install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/odind
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/yoda
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/kvasir

install-yoda: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/yoda

install-kvasir: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/kvasir

build: go.sum
	go build -mod=readonly -o ./build/odind $(BUILD_FLAGS) ./cmd/odind
	go build -mod=readonly -o ./build/yoda $(BUILD_FLAGS) ./cmd/yoda
	go build -mod=readonly -o ./build/kvasir $(BUILD_FLAGS) ./cmd/kvasir

build-yoda: go.sum
	go build -mod=readonly -o ./build/yoda $(BUILD_FLAGS) ./cmd/yoda

build-kvasir: go.sum
	go build -mod=readonly -o ./build/kvasir $(BUILD_FLAGS) ./cmd/kvasir

faucet: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/faucet

release: go.sum
	env GOOS=linux GOARCH=amd64 \
		go build -mod=readonly -o ./build/odind_linux_amd64 $(BUILD_FLAGS) ./cmd/odind
	env GOOS=darwin GOARCH=amd64 \
		go build -mod=readonly -o ./build/odind_darwin_amd64 $(BUILD_FLAGS) ./cmd/odind
	env GOOS=windows GOARCH=amd64 \
		go build -mod=readonly -o ./build/odind_windows_amd64 $(BUILD_FLAGS) ./cmd/odind
	env GOOS=linux GOARCH=amd64 \
		go build -mod=readonly -o ./build/yoda_linux_amd64 $(BUILD_FLAGS) ./cmd/yoda
	env GOOS=darwin GOARCH=amd64 \
		go build -mod=readonly -o ./build/yoda_darwin_amd64 $(BUILD_FLAGS) ./cmd/yoda
	env GOOS=windows GOARCH=amd64 \
		go build -mod=readonly -o ./build/yoda_windows_amd64 $(BUILD_FLAGS) ./cmd/yoda

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify
	touch go.sum

test:
	@go test -mod=readonly $(PACKAGES)

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.13.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@$(protoImage) sh ./scripts/protoc-swagger-gen.sh
	$(MAKE) update-swagger-docs

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-gen-any:
	$(DOCKER) run --rm -v $(pwd):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen-any.sh

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-check-breaking:
	@$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main

deb:
	rm -rf /tmp/GeoDB

	mkdir -p /tmp/GeoDB/deb/$(DEB_BIN_DIR)
	cp -f ./build/yoda /tmp/GeoDB/deb/$(DEB_BIN_DIR)/yoda
	cp -f ./build/odind /tmp/GeoDB/deb/$(DEB_BIN_DIR)/odind
	chmod +x /tmp/GeoDB/deb/$(DEB_BIN_DIR)/odind /tmp/GeoDB/deb/$(DEB_BIN_DIR)/yoda

	mkdir -p /tmp/GeoDB/deb/$(DEB_LIB_DIR)

	mkdir -p /tmp/GeoDB/deb/DEBIAN
	cp ./deployment/deb/control /tmp/GeoDB/deb/DEBIAN/control
	printf "Version: " >> /tmp/GeoDB/deb/DEBIAN/control
	printf "$(VERSION)" >> /tmp/GeoDB/deb/DEBIAN/control
	echo "" >> /tmp/GeoDB/deb/DEBIAN/control
	#cp ./deployment/deb/postinst /tmp/GeoDB/deb/DEBIAN/postinst
	#chmod 755 /tmp/GeoDB/deb/DEBIAN/postinst
	#cp ./deployment/deb/postrm /tmp/GeoDB/deb/DEBIAN/postrm
	#chmod 755 /tmp/GeoDB/deb/DEBIAN/postrm
	#cp ./deployment/deb/triggers /tmp/GeoDB/deb/DEBIAN/triggers
	#chmod 755 /tmp/GeoDB/deb/DEBIAN/triggers
	dpkg-deb --build /tmp/GeoDB/deb/ .
	-rm -rf /tmp/GeoDB
	cp ./odinprotocol_$(VERSION)_amd64.deb ./odinprotocol_v$(VERSION)_amd64.deb

.PHONY: proto-all proto-gen proto-swagger-gen proto-format proto-gen-any proto-lint proto-check-breaking
