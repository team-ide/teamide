@echo off

cd ../

go env -w CGO_ENABLED=1
go env -w GOOS=windows
go env -w GOARCH=amd64

go run . --isDev --passWindow

pause