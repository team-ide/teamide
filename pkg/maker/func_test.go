package maker

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/team-ide/go-tool/util"
	"testing"
)

func TestScript(t *testing.T) {
	run := goja.New()

	script := `var ss = 1;
(function (){
var a user.userId;
var aa []user;
ss++;
return ss
})()`

	p, err := goja.Parse("", script)
	if err != nil {
		panic(err)
	}

	fmt.Println(util.GetStringValue(p))
	res, err := run.RunScript("", script)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Export())
	compile, err := run.CompileScript("", `
(function (){
ss++;
return `+"`sss-${ss}`"+`
})()
`)
	if err != nil {
		panic(err)
	}
	res, err = run.RunProgram(compile)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Export())
	res, err = run.RunProgram(compile)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Export())
}
