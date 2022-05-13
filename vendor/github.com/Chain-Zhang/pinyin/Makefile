PKG_LIST := $(shell go list ./... | grep -v /vendor/)

test:
	@go test -short -cover ${PKG_LIST}