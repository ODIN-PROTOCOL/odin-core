PACKAGES=$(shell go list ./... | grep -v '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)

DEB_BIN_DIR ?= /usr/local/bin
DEB_LIB_DIR ?= /usr/lib

ifeq ($(LEDGER_ENABLED),true)
	build_tags += ledger
endif

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=odinchain \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=odind \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags)"

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/odind
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/yoda
	
install-yoda: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/yoda

build: go.sum
	go build -mod=readonly -o ./build/odind $(BUILD_FLAGS) ./cmd/odind
	go build -mod=readonly -o ./build/yoda $(BUILD_FLAGS) ./cmd/yoda

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

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen.sh

proto-format:
	@echo "Formatting Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace \
	--workdir /workspace tendermintdev/docker-build-proto \
	find ./ -not -path "./third_party/*" -name *.proto -exec clang-format -i {} \;

# This generates the SDK's custom wrapper for google.protobuf.Any. It should only be run manually when needed
proto-gen-any:
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen-any.sh

proto-swagger-gen:
	@./scripts/protoc-swagger-gen.sh

proto-lint:
	@$(DOCKER_BUF) check lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) check breaking --against-input $(HTTPS_GIT)#branch=master

TM_URL              = https://raw.githubusercontent.com/tendermint/tendermint/v0.34.0-rc6/proto/tendermint
GOGO_PROTO_URL      = https://raw.githubusercontent.com/regen-network/protobuf/cosmos
COSMOS_PROTO_URL    = https://raw.githubusercontent.com/regen-network/cosmos-proto/master
CONFIO_URL          = https://raw.githubusercontent.com/confio/ics23/v0.6.3

TM_CRYPTO_TYPES     = third_party/proto/tendermint/crypto
TM_ABCI_TYPES       = third_party/proto/tendermint/abci
TM_TYPES            = third_party/proto/tendermint/types
TM_VERSION          = third_party/proto/tendermint/version
TM_LIBS             = third_party/proto/tendermint/libs/bits
TM_P2P              = third_party/proto/tendermint/p2p

GOGO_PROTO_TYPES    = third_party/proto/gogoproto
COSMOS_PROTO_TYPES  = third_party/proto/cosmos_proto
CONFIO_TYPES        = third_party/proto/confio

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

proto-update-deps:
	@mkdir -p $(GOGO_PROTO_TYPES)
	@curl -sSL $(GOGO_PROTO_URL)/gogoproto/gogo.proto > $(GOGO_PROTO_TYPES)/gogo.proto

	@mkdir -p $(COSMOS_PROTO_TYPES)
	@curl -sSL $(COSMOS_PROTO_URL)/cosmos.proto > $(COSMOS_PROTO_TYPES)/cosmos.proto

## Importing of tendermint protobuf definitions currently requires the
## use of `sed` in order to build properly with cosmos-sdk's proto file layout
## (which is the standard Buf.build FILE_LAYOUT)
## Issue link: https://github.com/tendermint/tendermint/issues/5021
	@mkdir -p $(TM_ABCI_TYPES)
	@curl -sSL $(TM_URL)/abci/types.proto > $(TM_ABCI_TYPES)/types.proto

	@mkdir -p $(TM_VERSION)
	@curl -sSL $(TM_URL)/version/types.proto > $(TM_VERSION)/types.proto

	@mkdir -p $(TM_TYPES)
	@curl -sSL $(TM_URL)/types/types.proto > $(TM_TYPES)/types.proto
	@curl -sSL $(TM_URL)/types/evidence.proto > $(TM_TYPES)/evidence.proto
	@curl -sSL $(TM_URL)/types/params.proto > $(TM_TYPES)/params.proto
	@curl -sSL $(TM_URL)/types/validator.proto > $(TM_TYPES)/validator.proto
	@curl -sSL $(TM_URL)/types/block.proto > $(TM_TYPES)/block.proto

	@mkdir -p $(TM_CRYPTO_TYPES)
	@curl -sSL $(TM_URL)/crypto/proof.proto > $(TM_CRYPTO_TYPES)/proof.proto
	@curl -sSL $(TM_URL)/crypto/keys.proto > $(TM_CRYPTO_TYPES)/keys.proto

	@mkdir -p $(TM_LIBS)
	@curl -sSL $(TM_URL)/libs/bits/types.proto > $(TM_LIBS)/types.proto

	@mkdir -p $(TM_P2P)
	@curl -sSL $(TM_URL)/p2p/types.proto > $(TM_P2P)/types.proto

	@mkdir -p $(CONFIO_TYPES)
	@curl -sSL $(CONFIO_URL)/proofs.proto > $(CONFIO_TYPES)/proofs.proto
## insert go package option into proofs.proto file
## Issue link: https://github.com/confio/ics23/issues/32
	@sed -i '4ioption go_package = "github.com/confio/ics23/go";' $(CONFIO_TYPES)/proofs.proto

protoVer=v0.2
protoImageName=tendermintdev/sdk-proto-gen:$(protoVer)
containerProtoGen=$(PROJECT_NAME)-proto-gen-$(protoVer)
containerProtoGenAny=$(PROJECT_NAME)-proto-gen-any-$(protoVer)
containerProtoGenSwagger=$(PROJECT_NAME)-proto-gen-swagger-$(protoVer)
containerProtoFmt=$(PROJECT_NAME)-proto-fmt-$(protoVer)


proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGenSwagger}$$"; then docker start -a $(containerProtoGenSwagger); else docker run --name $(containerProtoGenSwagger) -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) \
		sh ./scripts/protoc-swagger-gen.sh; fi

.PHONY: proto-all proto-gen proto-gen-any proto-swagger-gen proto-format proto-lint proto-check-breaking proto-update-deps
