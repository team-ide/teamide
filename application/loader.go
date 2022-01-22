package application

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"teamide/application/base"
	"teamide/application/model"
)

type ContextLoader struct {
	Dir             string
	dirAbsolutePath string
}

func (this_ *ContextLoader) Load() (context *model.ModelContext, err error) {
	var path string
	path, err = os.Getwd()
	if err != nil {
		return
	}
	dirPath := path + "/" + this_.Dir
	var exists bool
	exists, err = base.PathExists(dirPath)
	if err != nil {
		return
	}
	if !exists {
		err = errors.New(fmt.Sprint("dir [" + dirPath + "] not exist"))
		return
	}
	var abs string
	abs, err = filepath.Abs(dirPath)
	if err != nil {
		return
	}
	this_.dirAbsolutePath = filepath.ToSlash(abs)

	fileContentMap := make(map[string]string)
	err = this_.loadDirFiles(fileContentMap, this_.dirAbsolutePath)
	var namePaths []string
	for namePath := range fileContentMap {
		namePaths = append(namePaths, namePath)
	}

	sort.Strings(namePaths)

	context = &model.ModelContext{}
	for _, namePath := range namePaths {
		if strings.HasSuffix(namePath, ".md") {
			continue
		}
		text := fileContentMap[namePath]
		if strings.Index(namePath, "struct/") == 0 {
			var one *model.StructModel
			one, err = model.TextToStructModel(strings.TrimPrefix(namePath, "struct/"), text)
			if err != nil {
				fmt.Println("to struct model error:", err)
				return
			}
			context.AppendStruct(one)
		} else if strings.Index(namePath, "service/") == 0 {
			var one *model.ServiceModel
			one, err = model.TextToServiceModel(strings.TrimPrefix(namePath, "service/"), text)
			if err != nil {
				fmt.Println("to service model error:", err)
				return
			}
			context.AppendService(one)
		} else if strings.Index(namePath, "test/") == 0 {
			var one *model.TestModel
			one, err = model.TextToTestModel(strings.TrimPrefix(namePath, "test/"), text)
			if err != nil {
				fmt.Println("to test model error:", err)
				return
			}
			context.AppendTest(one)
		} else if strings.Index(namePath, "teamide/server/web/") == 0 {
			var one *model.ServerWebModel
			one, err = model.TextToServerWebModel(strings.TrimPrefix(namePath, "teamide/server/web/"), text)
			if err != nil {
				fmt.Println("to server web model error:", err)
				return
			}
			context.AppendServerWeb(one)
		} else if strings.Index(namePath, "dictionary/") == 0 {
			var one *model.DictionaryModel
			one, err = model.TextToDictionaryModel(strings.TrimPrefix(namePath, "dictionary/"), text)
			if err != nil {
				fmt.Println("to dictionary model error:", err)
				return
			}
			context.AppendDictionary(one)
		} else if strings.Index(namePath, "constant/") == 0 {
			var one *model.ConstantModel
			one, err = model.TextToConstantModel(strings.TrimPrefix(namePath, "constant/"), text)
			if err != nil {
				fmt.Println("to constant model error:", err)
				return
			}
			context.AppendConstant(one)
		} else if strings.Index(namePath, "error/") == 0 {
			var one *model.ErrorModel
			one, err = model.TextToErrorModel(strings.TrimPrefix(namePath, "error/"), text)
			if err != nil {
				fmt.Println("to error model error:", err)
				return
			}
			context.AppendError(one)
		} else if strings.Index(namePath, "datasource/database/") == 0 {
			var one *model.DatasourceDatabase
			one, err = model.TextToDatasourceDatabase(strings.TrimPrefix(namePath, "datasource/database/"), text)
			if err != nil {
				fmt.Println("to datasource database model error:", err)
				return
			}
			context.AppendDatasourceDatabase(one)
		} else if strings.Index(namePath, "datasource/redis/") == 0 {
			var one *model.DatasourceRedis
			one, err = model.TextToDatasourceRedis(strings.TrimPrefix(namePath, "datasource/redis/"), text)
			if err != nil {
				fmt.Println("to datasource redis model error:", err)
				return
			}
			context.AppendDatasourceRedis(one)
		} else if strings.Index(namePath, "datasource/kafka/") == 0 {
			var one *model.DatasourceKafka
			one, err = model.TextToDatasourceKafka(strings.TrimPrefix(namePath, "datasource/kafka/"), text)
			if err != nil {
				fmt.Println("to datasource kafka model error:", err)
				return
			}
			context.AppendDatasourceKafka(one)
		} else if strings.Index(namePath, "datasource/zookeeper/") == 0 {
			var one *model.DatasourceZookeeper
			one, err = model.TextToDatasourceZookeeper(strings.TrimPrefix(namePath, "datasource/zookeeper/"), text)
			if err != nil {
				fmt.Println("to datasource zookeeper model error:", err)
				return
			}
			context.AppendDatasourceZookeeper(one)
		}
	}
	return
}

func (this_ *ContextLoader) loadDirFiles(fileContentMap map[string]string, fileDir string) (err error) {
	//获取当前目录下的所有文件或目录信息
	err = filepath.Walk(fileDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {

		} else {
			var abs string
			abs, err = filepath.Abs(path)
			if err != nil {
				return err
			}
			fileAbsolutePath := filepath.ToSlash(abs)
			name := strings.TrimPrefix(fileAbsolutePath, this_.dirAbsolutePath)
			name = strings.TrimPrefix(name, "/")
			var f *os.File
			f, err = os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			r := bufio.NewReader(f)
			var line string
			var content string
			for {
				line, err = r.ReadString('\n')
				if err != nil && err != io.EOF {
					return err
				}
				content += line
				if err == io.EOF {
					err = nil
					break
				}
			}
			fileContentMap[name] = content
		}
		return nil
	})
	return
}
