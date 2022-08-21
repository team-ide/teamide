@echo off

cd ../pkg/node/main

go run . -id node7 -address :21097 -token x -connAddress 127.0.0.1:21096 -connToken x

pause