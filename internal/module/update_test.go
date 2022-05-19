package module

import (
	"fmt"
	"testing"
)

func TestReleasesCheck(t *testing.T) {
	releaseHtml, err := releasesCheck()
	if err != nil {
		panic(err)
	}
	fmt.Println(len(releaseHtml))
	fmt.Println(releaseHtml)
}
