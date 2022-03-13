package static

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	BaseDir string
)

func init() {
	var err error
	BaseDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	BaseDir, err = filepath.Abs(BaseDir)
	if err != nil {
		panic(err)
	}
	BaseDir = filepath.ToSlash(BaseDir)
	if !strings.HasSuffix(BaseDir, "/") {
		BaseDir += "/"
	}
}

// test -timeout 3600s -run ^TestStatic$ teamide/static
func TestStatic(t *testing.T) {

	var err error
	var dist string
	dist, err = filepath.Abs(BaseDir + "../html/dist")
	if err != nil {
		panic(err)
	}
	dist = filepath.ToSlash(dist)
	err = SetAsset(dist, BaseDir+"html.go")
	if err != nil {
		panic(err)
	}
	var bs []byte = Asset("index.html")
	fmt.Println(string(bs))

}
