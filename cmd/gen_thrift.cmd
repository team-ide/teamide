@echo on

cd %~dp0
cd ../

thrift -out ./pkg --gen "go:package_prefix=teamide/pkg/,thrift_import=github.com/apache/thrift/lib/go/thrift" thrift/bean.thrift
