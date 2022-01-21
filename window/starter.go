package window

import (
	"net/url"
	"runtime"

	"github.com/zserge/lorca"
)

var (
	window_width  = getWindowWidth()
	window_height = getWindowHeight()
)

func init() {
	if window_width == 0 {
		window_width = 1024
	}
	if window_height == 0 {
		window_height = 768
	}
	window_width = window_width - 20
	window_height = window_height - 60
}

func Start(title string, webUrl string, onClose func()) (err error) {

	err = startWindow(title, webUrl, onClose)

	if err != nil {
		return
	}

	return
}

func startWindow(title string, webUrl string, onClose func()) (err error) {
	var args []string
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
		<title>`+title+`</title>
		<link rel="icon" href="`+webUrl+`static/favicon.png">
		<script type="text/javascript">
			location.href="`+webUrl+`"
		</script>
	</head>
</html>`), "", window_width, window_height, args...)
	if err != nil {
		return
	}
	go func() {
		if onClose != nil {
			defer onClose()
		}
		defer ui.Close()
		<-ui.Done()

	}()
	return
}
