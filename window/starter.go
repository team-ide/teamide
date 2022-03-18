package window

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
	window_width = window_width - 40
	window_height = window_height - 80
}

func Start(title string, webUrl string, onClose func()) (err error) {

	err = startWindow(title, webUrl, onClose)

	if err != nil {
		return
	}

	return
}

var defaultChromeArgs = []string{
	"--disable-background-networking",
	"--disable-background-timer-throttling",
	"--disable-backgrounding-occluded-windows",
	"--disable-breakpad",
	"--disable-client-side-phishing-detection",
	"--disable-default-apps",
	"--disable-dev-shm-usage",
	"--disable-infobars",
	"--disable-extensions",
	"--disable-features=site-per-process",
	"--disable-hang-monitor",
	"--disable-ipc-flooding-protection",
	"--disable-popup-blocking",
	"--disable-prompt-on-repost",
	"--disable-renderer-backgrounding",
	"--disable-sync",
	"--disable-translate",
	"--disable-windows10-custom-titlebar",
	"--metrics-recording-only",
	"--no-first-run",
	"--no-default-browser-check",
	"--safebrowsing-disable-auto-update",
	"--enable-automation",
	"--password-store=basic",
	"--use-mock-keychain",
}

func startWindow(title string, webUrl string, onClose func()) (err error) {

	locateChrome := LocateChrome()
	if locateChrome != "" {
		args := []string{}
		args = append(args, defaultChromeArgs...)
		args = append(args, "--app="+webUrl)
		args = append(args, "--class=TeamIDE")
		args = append(args, "--start-maximized")
		args = append(args, "--window-position=20,20") // 窗口位置
		args = append(args, fmt.Sprintf("--window-size=%d,%d", window_width, window_height))
		cmd := exec.Command(locateChrome, args...)

		err = cmd.Run()

		if err != nil {
			return
		}
		return
	}

	// 	var args []string
	// 	args = append(args, "--class=TeamIDE")
	// 	// args = append(args, "--kiosk")               // 最大化
	// 	// args = append(args, "--start-maximized")     // 最大化
	// 	args = append(args, "--window-position=20,20") // 窗口位置

	// 	// args = append(args, "--window-size=-1,-1")   // 强制显示更新菜单项
	// 	var ui lorca.UI
	// 	ui, err = lorca.New("data:text/html,"+url.PathEscape(`
	// <html>
	// 	<head>
	// 		<meta charset="utf-8">
	// 		<meta http-equiv="X-UA-Compatible" content="IE=edge">
	// 		<meta name="viewport" content="width=device-width,initial-scale=1.0">
	// 		<title>`+title+`</title>
	// 		<link rel="icon" href="`+webUrl+`static/favicon.png">
	// 		<script type="text/javascript">
	// 			location.href="`+webUrl+`"
	// 		</script>
	// 	</head>
	// </html>`), "", window_width, window_height, args...)
	// 	if err != nil {
	// 		return
	// 	}
	// 	go func() {
	// 		if onClose != nil {
	// 			defer onClose()
	// 		}
	// 		defer ui.Close()
	// 		<-ui.Done()

	// 	}()
	return
}

func LocateChrome() string {

	// If env variable "LORCACHROME" specified and it exists
	if path, ok := os.LookupEnv("LORCACHROME"); ok {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	var paths []string
	switch runtime.GOOS {
	case "darwin":
		paths = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/google-chrome",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
		}
	case "windows":
		paths = []string{
			os.Getenv("LocalAppData") + "/Google/Chrome/Application/chrome.exe",
			os.Getenv("ProgramFiles") + "/Google/Chrome/Application/chrome.exe",
			os.Getenv("ProgramFiles(x86)") + "/Google/Chrome/Application/chrome.exe",
			os.Getenv("LocalAppData") + "/Chromium/Application/chrome.exe",
			os.Getenv("ProgramFiles") + "/Chromium/Application/chrome.exe",
			os.Getenv("ProgramFiles(x86)") + "/Chromium/Application/chrome.exe",
			os.Getenv("ProgramFiles(x86)") + "/Microsoft/Edge/Application/msedge.exe",
		}
	default:
		paths = []string{
			"/usr/bin/google-chrome-stable",
			"/usr/bin/google-chrome",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/snap/bin/chromium",
		}
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		return path
	}
	return ""
}
