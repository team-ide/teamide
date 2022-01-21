package main

import (
	"net/url"
	"runtime"
	"syscall"

	"github.com/zserge/lorca"
)

var (
	window_width  = GetSystemMetrics(0) - 20
	window_height = GetSystemMetrics(1) - 60
)

func GetSystemMetrics(nIndex int) int {
	ret, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(uintptr(nIndex))
	return int(ret)
}
func openWindow(webServerAddress string) (err error) {
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	// args = append(args, "--single-proces")   // 单进程运行
	// args = append(args, "--start-maximized") // 启动就最大化
	// args = append(args, "--no-sandbox")      // 取消沙盒模式
	var ui lorca.UI
	ui, err = lorca.New("data:text/html,"+url.PathEscape(`
<html>
	<head>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width,initial-scale=1.0">
		<title>Team IDE · Model Coder</title>
		<link rel="icon" href="`+webServerAddress+`static/favicon.png">
		<script type="text/javascript">
			location.href="`+webServerAddress+`"
		</script>
	</head>
</html>
	`), "", window_width, window_height, args...)
	if err != nil {
		return
	}
	// err = ui.Load(webServerAddress)
	// if err != nil {
	// 	ui.Close()
	// 	return
	// }
	go func() {
		defer ui.Close()
		<-ui.Done()
		onWindowClose()
	}()
	// time.Sleep(1000 * 3)
	return
}
