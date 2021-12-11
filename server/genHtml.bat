
cd ../html/dist

go-bindata -pkg web -o ../../server/web/html.go ./...

pause