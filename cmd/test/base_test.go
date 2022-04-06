package test

import (
	"fmt"
	"testing"
)

func TestIDMysql(t *testing.T) {
	//var err error
	str := "teamide:event:"
	fmt.Println(str)
	bs := []byte(str)
	fmt.Println(len(bs))

}
