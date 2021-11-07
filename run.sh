#!/bin/bash
mkdir ./bin

cd ./migrations/
tern migrate -m .
tern code install .

cd ../cmd/blockchain-info-server/
go build -o ../../bin/bc_info .

cd ../..
./bin/bc_info --port 8080 --host 0.0.0.0