ARG ONOS_E2T_BASE_VERSION=latest

FROM onosproject/onos-e2sub-base:$ONOS_E2T_BASE_VERSION as base

FROM alpine:3.11
RUN apk add libc6-compat

USER nobody

COPY --from=base /go/src/github.com/onosproject/onos-e2sub/build/_output/onos-e2sub /usr/local/bin/onos-e2sub

ENTRYPOINT ["onos-e2sub"]
