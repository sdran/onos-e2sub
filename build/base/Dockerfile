ARG ONOS_BUILD_VERSION=undefined

FROM onosproject/golang-build:$ONOS_BUILD_VERSION

ENV GO111MODULE=on
ARG ONOS_MAKE_TARGET=build

COPY Makefile go.mod go.sum /go/src/github.com/onosproject/onos-e2sub/
COPY cmd/ /go/src/github.com/onosproject/onos-e2sub/cmd/
COPY pkg/ /go/src/github.com/onosproject/onos-e2sub/pkg/
COPY vendor/ /go/src/github.com/onosproject/onos-e2sub/vendor/
COPY api/ /go/src/github.com/onosproject/onos-e2sub/api

RUN cd /go/src/github.com/onosproject/onos-e2sub && GOFLAGS=-mod=vendor make ${ONOS_MAKE_TARGET}
