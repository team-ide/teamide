package static

import (
	"fmt"
	"path/filepath"
	"teamide/util"
	"testing"
)

// test -timeout 3600s -run ^TestStatic$ teamide/static
func TestStatic(t *testing.T) {

	var err error
	var dist string
	dist, err = filepath.Abs(util.BaseDir + "../html/dist")
	if err != nil {
		panic(err)
	}
	dist = filepath.ToSlash(dist)
	err = SetAsset(dist, util.BaseDir+"html.go")
	if err != nil {
		panic(err)
	}
	var bs []byte = Asset("index.html")
	fmt.Println(string(bs))

}
