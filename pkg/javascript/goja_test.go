package javascript

import (
	"fmt"
	"testing"
)

func TestAddTask(t *testing.T) {
	script := "1 + _$uuid()"
	context := GetContext()

	res, err := Run(script, context)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
