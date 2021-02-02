FROM onosproject/golang-build:v0.6.7 as build

ENV GO111MODULE=on
ARG ONOS_MAKE_TARGET=build

COPY Makefile go.mod go.sum /go/src/github.com/onosproject/onos-e2sub/
COPY cmd/ /go/src/github.com/onosproject/onos-e2sub/cmd/
COPY pkg/ /go/src/github.com/onosproject/onos-e2sub/pkg/
COPY vendor/ /go/src/github.com/onosproject/onos-e2sub/vendor/

RUN cd /go/src/github.com/onosproject/onos-e2sub && GOFLAGS=-mod=vendor make ${ONOS_MAKE_TARGET}

FROM alpine:3.11
RUN apk add libc6-compat

USER nobody

COPY --from=build /go/src/github.com/onosproject/onos-e2sub/build/_output/onos-e2sub /usr/local/bin/onos-e2sub

ENTRYPOINT ["onos-e2sub"]
