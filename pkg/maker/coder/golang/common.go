package golang

import (
	"sort"
	"strings"
)

var (
	commonCode = `
var (
	ReleaseTime    = ""
	ReleaseVersion = "0.0.1"
	GitCommit      = "xxx"
	RootDir        = ""
	UserHomeDir    = ""
)

type Error struct {
	code string
	msg  string
}

func (this_ *Error) Error() string {
	return fmt.Sprintf("code:%s , msg:%s", this_.code, this_.msg)
}

// NewError 构造异常对象，code为错误码，msg为错误信息
func NewError(code string, msg string) *Error {
	err := &Error{
		code: code,
		msg:  msg,
	}
	return err
}
`
)

func (this_ *Generator) GenCommon() (err error) {
	dir := this_.golang.GetCommonDir(this_.Dir)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "common.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	var imports []string

	pack := this_.golang.GetCommonPack()

	builder.AppendTabLine("package " + pack)
	builder.NewLine()

	builder.AppendTabLine("import(")
	builder.Tab()

	ss := strings.Split(`
	"fmt"
`, "\n")
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		s = strings.TrimPrefix(s, `"`)
		s = strings.TrimSuffix(s, `"`)
		imports = append(imports, s)
	}

	sort.Strings(imports)
	for _, im := range imports {
		as := strings.Split(im, " ")
		if len(as) == 2 {
			builder.AppendTabLine(as[1] + " \"" + as[0] + "\"")
		} else {
			builder.AppendTabLine("\"" + im + "\"")
		}
	}
	builder.Indent()
	builder.AppendTabLine(")")
	builder.NewLine()

	code := strings.ReplaceAll(commonCode, "{pack}", pack)

	builder.AppendCode(code)
	err = this_.GenCommonEvent()
	return
}

var (
	commonEventCode = `
type Event string

const (
	EventAppStart = Event("app-start-event")

	EventAppConfigBefore = Event("app-config-before-event")
	EventAppConfigAfter  = Event("app-config-after-event")

	EventAppComponentBefore = Event("app-component-before-event")
	EventAppComponentAfter  = Event("app-component-after-event")

	EventAppIFaceBefore = Event("app-iFace-before-event")
	EventAppIFaceAfter  = Event("app-iFace-after-event")

	EventAppReady = Event("app-ready-event")

	EventAppStop = Event("app-stop-event")
)

var (
	eventListeners     = map[Event]listenerArray{}
	eventListenersLock sync.Mutex
)

type Listener struct {
	event   Event
	order   int
	onEvent func(args ...any)
}

func OnEvent(event Event, onEvent func(args ...any), order int) {
	eventListenersLock.Lock()
	defer eventListenersLock.Unlock()

	listener := &Listener{
		event:   event,
		onEvent: onEvent,
		order:   order,
	}
	eventListeners[event] = append(eventListeners[event], listener)
	eventListeners[event].Sort()
}

func GetListeners(event Event) (res []*Listener) {
	eventListenersLock.Lock()
	defer eventListenersLock.Unlock()

	res = eventListeners[event]
	return
}

func CallEvent(event Event, args ...any) {
	listeners := GetListeners(event)
	if len(listeners) == 0 {
		return
	}
	for _, listener := range listeners {
		doEvent(listener, args...)
	}
	return
}

func doEvent(listener *Listener, args ...any) {
	if listener == nil || listener.onEvent == nil {
		return
	}
	defer func() {
		if e := recover(); e != nil {
			err := errors.New(fmt.Sprint(e))
			fmt.Println("on event ["+listener.event+"] error:", err)
			logger.Logger.Error("doEvent", zap.Error(err))
		}
	}()

	listener.onEvent(args...)

	return
}

type listenerArray []*Listener

func (p listenerArray) Len() int           { return len(p) }
func (p listenerArray) Less(i, j int) bool { return p[i].order < p[j].order }
func (p listenerArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p listenerArray) Sort()              { sort.Sort(p) }

func OnSignal() {
	c := make(chan os.Signal)
	signal.Notify(c)
	for s := range c {
		switch s {
		case os.Kill: // kill -9 pid，下面的无效
			fmt.Println("强制退出", s)
			CallEvent(EventAppStop)
			os.Exit(0)
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT: // ctrl + c
			fmt.Println("退出", s)
			CallEvent(EventAppStop)
			os.Exit(0)
		}
	}
}

var (
	waitGroupForStop       *sync.WaitGroup
	waitGroupForStopLocker sync.Mutex
)

func Wait() {
	waitGroupForStopLocker.Lock()
	if waitGroupForStop != nil {
		waitGroupForStop = &sync.WaitGroup{}
		waitGroupForStop.Add(1)
	}
	waitGroupForStopLocker.Unlock()
	waitGroupForStop.Wait()
}

func CallStop() {
	waitGroupForStopLocker.Lock()
	if waitGroupForStop != nil {
		waitGroupForStop.Done()
		waitGroupForStop = nil
	}
	waitGroupForStopLocker.Unlock()
	os.Exit(0)
}
`
)

func (this_ *Generator) GenCommonEvent() (err error) {
	dir := this_.golang.GetCommonDir(this_.Dir)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "event.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	var imports []string

	pack := this_.golang.GetCommonPack()

	builder.AppendTabLine("package " + pack)
	builder.NewLine()

	builder.AppendTabLine("import(")
	builder.Tab()

	imports = append(imports, this_.golang.GetLoggerImport())

	ss := strings.Split(`
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
`, "\n")
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		s = strings.TrimPrefix(s, `"`)
		s = strings.TrimSuffix(s, `"`)
		imports = append(imports, s)
	}

	sort.Strings(imports)
	for _, im := range imports {
		as := strings.Split(im, " ")
		if len(as) == 2 {
			builder.AppendTabLine(as[1] + " \"" + as[0] + "\"")
		} else {
			builder.AppendTabLine("\"" + im + "\"")
		}
	}
	builder.Indent()
	builder.AppendTabLine(")")
	builder.NewLine()

	code := strings.ReplaceAll(commonEventCode, "{pack}", pack)

	builder.AppendCode(code)
	return
}
