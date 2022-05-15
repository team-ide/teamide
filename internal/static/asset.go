package static

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"teamide/pkg/util"
)

var (
	staticCache = map[string][]byte{}
)

func Asset(name string) []byte {
	bs, ok := staticCache[name]
	if ok {
		unzipBS, err := util.UGZipBytes(bs)
		if err != nil {
			util.Logger.Error("Asset["+name+"]异常", zap.Error(err))
			return nil
		}
		return unzipBS
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
		zipBS, err := util.GZipBytes(bs)
		if err != nil {
			util.Logger.Error("SetAsset["+filename+"]异常", zap.Error(err))
			return err
		}
		fmt.Println("文件[" + filename + "]大小[" + fmt.Sprint(len(bs)) + "]压缩后大小[" + fmt.Sprint(len(zipBS)) + "]")

		f.WriteString(`	staticCache["` + filename + `"] = ` + "[]byte{")
		size := len(zipBS)
		for i, b := range zipBS {
			if i == size-1 {
				f.WriteString(fmt.Sprintf("%d", b))
			} else {
				f.WriteString(fmt.Sprintf("%d,", b))
			}
		}
		f.WriteString("}")
		f.WriteString("\n")
		f.WriteString("\n")
		f.WriteString("\n")
	}
	f.WriteString("}" + "\n")
	return
}
