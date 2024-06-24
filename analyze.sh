#!/bin/bash
set -e

./campus --analyze

cp cache/GeneratedProto/*.proto proto/

protoc -I=./proto --go_out=. proto/octodb.proto
protoc -I=./proto --go_out=. --go_opt=module=vertesan/campus \
  proto/penum.proto proto/pcommon.proto proto/pmaster.proto proto/papicommon.proto proto/ptransaction.proto proto/papi.proto
protoc -I=./proto --go-grpc_out=. --go-grpc_opt=module=vertesan/campus proto/papi.proto

cp cache/GeneratedProto/mapping.go proto/mapping/
