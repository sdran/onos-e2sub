export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: build

ONOS_E2T_VERSION := latest
ONOS_BUILD_VERSION := v0.6.6
ONOS_PROTOC_VERSION := v0.6.6
BUF_VERSION := 0.27.1

build: # @HELP build the Go binaries and run all validations (default)
build: build-tools
	go build -o build/_output/onos-e2sub ./cmd/onos-e2sub

test: # @HELP run the unit tests and source code validation
test: build deps linters license_check
	go test -race github.com/onosproject/onos-e2sub/pkg/...
	go test -race github.com/onosproject/onos-e2sub/cmd/...

jenkins-test: # @HELP run the unit tests and source code validation producing a junit style report for Jenkins
jenkins-test: build deps linters license_check
	TEST_PACKAGES=github.com/onosproject/onos-e2sub/... ./../build-tools/build/jenkins/make-unit

coverage: # @HELP generate unit test coverage data
coverage: build deps linters license_check
	export GOPRIVATE="github.com/onosproject/*"
	go test -covermode=count -coverprofile=onos.coverprofile github.com/onosproject/onos-e2sub/pkg/...
	cd .. && go get github.com/mattn/goveralls && cd onos-e2sub
	grep -v .pb.go onos.coverprofile >onos-nogrpc.coverprofile
	goveralls -coverprofile=onos-nogrpc.coverprofile -service travis-pro -repotoken xZuVup4oLZFkqxtkFW2qEkFTf9NDZhN2g

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: golang-ci # @HELP examines Go source code and reports coding problems
	golangci-lint run --timeout 5m

build-tools: # @HELP install the ONOS build tools if needed
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi

jenkins-tools: # @HELP installs tooling needed for Jenkins
	cd .. && go get -u github.com/jstemmer/go-junit-report && go get github.com/t-yuki/gocover-cobertura

golang-ci: # @HELP install golang-ci if not present
	golangci-lint --version || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b `go env GOPATH`/bin v1.36.0

license_check: build-tools # @HELP examine and ensure license headers exist
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR} --boilerplate LicenseRef-ONF-Member-1.0

gofmt: # @HELP run the Go format validation
	bash -c "diff -u <(echo -n) <(gofmt -d pkg/ cmd/ tests/)"

buflint: #@HELP run the "buf check lint" command on the proto files in 'api'
	docker run -it -v `pwd`:/go/src/github.com/onosproject/onos-e2sub \
		-w /go/src/github.com/onosproject/onos-e2sub/api \
		bufbuild/buf:${BUF_VERSION} check lint

protos: # @HELP compile the protobuf files (using protoc-go Docker)
protos:
	docker run -it -v `pwd`:/go/src/github.com/onosproject/onos-e2sub \
		-w /go/src/github.com/onosproject/onos-e2sub \
		--entrypoint build/bin/compile-protos.sh \
		onosproject/protoc-go:${ONOS_PROTOC_VERSION}

onos-e2sub-docker: # @HELP build onos-e2sub Docker image
onos-e2sub-docker:
	docker build . -f build/onos-e2sub/Dockerfile \
		-t onosproject/onos-e2sub:${ONOS_E2T_VERSION}

images: # @HELP build all Docker images
images: build onos-e2sub-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-e2sub:${ONOS_E2T_VERSION}

all: build images

publish: # @HELP publish version on github and dockerhub
	./../build-tools/publish-version ${VERSION} onosproject/onos-e2sub

jenkins-publish: build-tools jenkins-tools # @HELP Jenkins calls this to publish artifacts
	./build/bin/push-images
	../build-tools/release-merge-commit

bumponosdeps: # @HELP update "onosproject" go dependencies and push patch to git. Add a version to dependency to make it different to $VERSION
	./../build-tools/bump-onos-deps ${VERSION}

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./cmd/onos-e2sub/onos-e2sub ./cmd/onos/onos
	go clean -testcache github.com/onosproject/onos-e2sub/...

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
