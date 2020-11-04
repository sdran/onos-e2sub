#!/bin/sh

proto_imports=".:${GOPATH}/src/github.com/gogo/protobuf/protobuf:${GOPATH}/src/github.com/gogo/protobuf:${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate:${GOPATH}/src":"${GOPATH}/src/github.com/onosproject/onos-e2sub/api"

protoc -I=$proto_imports  --doc_out=docs/api  --doc_opt=markdown,subscription.md  --gogofaster_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=github.com/onosproject/onos-e2sub/api/e2/registry/v1beta1,plugins=grpc:. api/e2/registry/v1beta1/registry.proto
