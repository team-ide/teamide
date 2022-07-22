@echo off

cd ../node/main

go run . -id node5 -address :21095 -token x -connAddress 127.0.0.1:21090 -connToken x

pause