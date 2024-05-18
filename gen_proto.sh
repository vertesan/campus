protoc -I=./proto --go_out=. proto/octodb.proto
protoc -I=./proto --go_out=. proto/mastertag.proto
protoc -I=./proto --go_out=. --go_opt=module=vertesan/campus \
  proto/penum.proto proto/pcommon.proto proto/pmaster.proto
