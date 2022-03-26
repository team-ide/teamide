package base

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	RootDir      string
	UserHomeDir  string
	IsStandAlone = false
	IsHtmlDev    = false
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
	RootDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	RootDir, err = filepath.Abs(RootDir)
	if err != nil {
		panic(err)
	}
	RootDir = filepath.ToSlash(RootDir)
	if !strings.HasSuffix(RootDir, "/") {
		RootDir += "/"
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
