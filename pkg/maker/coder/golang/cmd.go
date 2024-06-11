package golang

var (
	cmdBuildCmdCode = `@echo off
cd ../

go env -w CGO_ENABLED=1
go env -w GOOS=windows
go env -w GOARCH=amd64

go mod tidy

go build -ldflags "-s -X 'app/common.ReleaseVersion=0.0.1' -X 'app/common.ReleaseTime=2024-06-11 16:00' -X 'app/common.GitCommit=xxx'" -o app.exe .
`
)

func (this_ *Generator) GenCmd() (err error) {
	dir := this_.Dir + "cmd/"
	if err = this_.Mkdir(dir); err != nil {
		return
	}

	err = this_.GenCmdBuildCmd(dir)
	if err != nil {
		return
	}
	return
}

func (this_ *Generator) GenCmdBuildCmd(dir string) (err error) {
	path := dir + "build.cmd"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	code := cmdBuildCmdCode
	builder.AppendCode(code)
	return
}
