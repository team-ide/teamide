package static

import (
	"encoding/base64"
	"os"
	"teamide/util"
)

var (
	static_cache = map[string]string{}
)

func Asset(name string) []byte {
	var content string = static_cache[name]
	if content != "" {
		bs, _ := base64.StdEncoding.DecodeString(content)
		return bs
	}
	return nil
}

func SetAsset(dir string, saveFile string) (err error) {
	var fileMap map[string][]byte = map[string][]byte{}
	err = util.LoadDirFiles(fileMap, dir)
	if err != nil {
		return
	}
	var exists bool
	exists, err = util.PathExists(saveFile)
	if err != nil {
		return
	}
	if exists {
		err = os.Remove(saveFile)
		if err != nil {
			return
		}
	}
	var f *os.File
	f, err = os.Create(saveFile)
	if err != nil {
		return
	}
	defer f.Close()

	f.WriteString("package static" + "\n")
	f.WriteString("\n")
	f.WriteString("\n")
	f.WriteString("func init() {" + "\n")
	for filename, bs := range fileMap {
		f.WriteString(`	static_cache["` + filename + `"] = ` + "`")
		f.WriteString(base64.StdEncoding.EncodeToString(bs))
		f.WriteString("`")
		f.WriteString("\n")
		f.WriteString("\n")
		f.WriteString("\n")
	}
	f.WriteString("}" + "\n")
	return
}
