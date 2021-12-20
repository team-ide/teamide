package main

import (
	"fmt"
	"server/modelcoder"
	"testing"
)

func TestApplication(t *testing.T) {
	applicationModel := &modelcoder.ApplicationModel{}

	applicationModel.AppendConstant(&modelcoder.ConstantModel{Name: "TEST_1"})

	application := modelcoder.NewApplication(applicationModel)

	res, err := application.InvokeServiceByName("aa", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(ToJSON(res))
}
