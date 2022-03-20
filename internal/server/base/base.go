package base

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	BaseDir      string
	UserHomeDir  string
	IsStandAlone = false
)

func getUserHome() string {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir
	}
	return ""
}
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

	userHome := getUserHome()
	if userHome != "" {
		userHome, err = filepath.Abs(userHome)
		if err != nil {
			panic(err)
		}
		UserHomeDir = filepath.ToSlash(userHome)
		if !strings.HasSuffix(UserHomeDir, "/") {
			UserHomeDir += "/"
		}

	}

}
