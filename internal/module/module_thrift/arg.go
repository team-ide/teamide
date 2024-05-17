package module_thrift

import (
	"errors"
	"github.com/dop251/goja"
	"github.com/team-ide/go-tool/javascript"
	"github.com/team-ide/go-tool/task"
	"github.com/team-ide/go-tool/util"
	"regexp"
	"sync"
)

func newArgFormat() (res *argFormat, err error) {
	res = &argFormat{}
	res.runtime = goja.New()
	res.scriptContext = javascript.NewContext()
	if len(res.scriptContext) > 0 {
		for key, value := range res.scriptContext {
			err = res.runtime.Set(key, value)
			if err != nil {
				return
			}
		}
	}
	return
}

type argFormat struct {
	runtime       *goja.Runtime
	scriptContext map[string]interface{}
	lock          sync.Mutex
}

func (this_ *argFormat) scriptValue(script string, param *task.ExecutorParam) (res string, err error) {
	if script == "" {
		return
	}

	this_.lock.Lock()
	defer this_.lock.Unlock()

	if param == nil {
		param = &task.ExecutorParam{}
	}
	err = this_.runtime.Set("index", param.Index)
	if err != nil {
		return
	}
	err = this_.runtime.Set("workerIndex", param.WorkerIndex)
	if err != nil {
		return
	}

	v, err := this_.runtime.RunString(script)
	if err != nil {
		err = errors.New("get scriptValue error:" + err.Error())
		return
	}
	res = util.GetStringValue(v)

	return
}
func (this_ *argFormat) stringArg(arg string, param *task.ExecutorParam) (res interface{}, err error) {
	if arg == "" {
		res = ""
		return
	}
	text := ""
	var re *regexp.Regexp
	re, _ = regexp.Compile(`[$]+{(.+?)}`)
	indexList := re.FindAllIndex([]byte(arg), -1)
	var lastIndex int = 0
	for _, indexes := range indexList {
		text += arg[lastIndex:indexes[0]]

		lastIndex = indexes[1]

		script := arg[indexes[0]+2 : indexes[1]-1]
		v := ""
		v, err = this_.scriptValue(script, param)
		if err != nil {
			return
		}
		text += v
	}
	text += arg[lastIndex:]

	res = text
	return
}
func (this_ *argFormat) formatArg(arg interface{}, param *task.ExecutorParam) (res interface{}, err error) {
	if arg == nil {
		return
	}
	switch tV := arg.(type) {
	case string:
		res, err = this_.stringArg(tV, param)
		break
	case []interface{}:
		var list []interface{}
		for _, one := range tV {
			var v interface{}
			v, err = this_.formatArg(one, param)
			if err != nil {
				return
			}
			list = append(list, v)
		}
		res = list
		break
	case map[string]interface{}:
		var data = map[string]interface{}{}
		for key, one := range tV {
			var v interface{}
			v, err = this_.formatArg(one, param)
			if err != nil {
				return
			}
			data[key] = v
		}
		res = data
		break
	default:
		res = tV
		break
	}

	return
}

func (this_ *argFormat) formatArgs(args []interface{}, param *task.ExecutorParam) (res []interface{}, err error) {
	if len(args) == 0 {
		return
	}
	for _, arg := range args {
		var v interface{}
		v, err = this_.formatArg(arg, param)
		if err != nil {
			return
		}
		res = append(res, v)
	}

	return
}
