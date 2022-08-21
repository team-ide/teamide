@echo off

cd ../pkg/node/main

go run . -id node1 -address :21091 -token x -connAddress 127.0.0.1:21090 -connToken x

pause