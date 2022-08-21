@echo off

cd ../pkg/node/main

go run . -id node2 -address :21092 -token x -connAddress 127.0.0.1:21090 -connToken x

pause