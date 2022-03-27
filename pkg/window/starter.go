package window

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
)

var (
	windowWidth  = getWindowWidth()
	windowHeight = getWindowHeight()
)

func init() {
	if windowWidth == 0 {
		windowWidth = 1440
	}
	if windowHeight == 0 {
		windowHeight = 900
	}
}

func Start(webUrl string, onClose func()) (err error) {

	err = startWindow(webUrl, onClose)

	if err != nil {
		return
	}

	return
}

func startWindow(webUrl string, onClose func()) (err error) {

	locateChrome := LocateChrome()
	if locateChrome == "" {
		err = errors.New("未检测到谷歌浏览器")
		return
	}
	var tmpDir string
	tmpDir, err = ioutil.TempDir("", "TeamIDE")
	if err != nil {
		return
	}
	var listener net.Listener
	listener, err = net.Listen("tcp", ":0")
	if err != nil {
		return
	}
	debuggingPort := listener.Addr().(*net.TCPAddr).Port

	var args []string
	//args = append(args, defaultChromeArgs...)
	args = append(args, "--app="+webUrl)
	args = append(args, fmt.Sprintf("--user-data-dir=%s", tmpDir))
	args = append(args, "--class=TeamIDE")
	//args = append(args, "--start-maximized")

	width := 1440
	height := 900
	left := (windowWidth - width) / 2
	right := (windowHeight - height) / 2

	args = append(args, fmt.Sprintf("--window-position=%d,%d", left, right)) // 窗口位置
	args = append(args, fmt.Sprintf("--window-size=%d,%d", width, height))
	args = append(args, fmt.Sprintf("--remote-debugging-port=%d", debuggingPort))

	b := &browser{
		name: locateChrome,
		args: args,
	}
	b.cmd = exec.Command(locateChrome, args...)

	var pipe io.ReadCloser
	if pipe, err = b.cmd.StderrPipe(); err != nil {
		return
	}

	if err = b.cmd.Start(); err != nil {
		return
	}

	re := regexp.MustCompile(`^DevTools listening on (ws://.*?)\r?\n$`)
	var m []string
	m, err = readUntilMatch(pipe, re)
	if err != nil {
		b.kill()
		return
	}
	wsURL := m[1]

	//fmt.Println("wsURL:", wsURL)
	// Open a websocket
	b.ws, err = websocket.Dial(wsURL, "", "http://127.0.0.1")
	if err != nil {
		b.kill()
		return
	}

	go func() {
		err = b.cmd.Wait()
		if err != nil {
			b.kill()
		}
		onClose()
	}()
	return

	//var args []string
	//args = append(args, "--class=TeamIDE")
	//// args = append(args, "--kiosk")               // 最大化
	//// args = append(args, "--start-maximized")     // 最大化
	//args = append(args, "--window-position=20,20") // 窗口位置
	//
	//// args = append(args, "--window-size=-1,-1")   // 强制显示更新菜单项
	//var ui lorca.UI
	//ui, err = lorca.New("data:text/html,"+url.PathEscape(`
	//<html>
	//	<head>
	//		<meta charset="utf-8">
	//		<meta http-equiv="X-UA-Compatible" content="IE=edge">
	//		<meta name="viewport" content="width=device-width,initial-scale=1.0">
	//		<title>`+title+`</title>
	//		<link rel="icon" href="`+webUrl+`static/favicon.png">
	//		<script type="text/javascript">
	//			location.href="`+webUrl+`"
	//		</script>
	//	</head>
	//</html>`), "", window_width, window_height, args...)
	//if err != nil {
	//	return
	//}
	//go func() {
	//	if onClose != nil {
	//		defer onClose()
	//	}
	//	defer ui.Close()
	//	<-ui.Done()
	//
	//}()
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

type browser struct {
	name string
	args []string
	cmd  *exec.Cmd
	ws   *websocket.Conn
}

func (c *browser) kill() error {
	if c.ws != nil {
		if err := c.ws.Close(); err != nil {
			return err
		}
	}
	if state := c.cmd.ProcessState; state == nil || !state.Exited() {
		return c.cmd.Process.Kill()
	}
	return nil
}

func readUntilMatch(r io.ReadCloser, re *regexp.Regexp) ([]string, error) {
	br := bufio.NewReader(r)
	for {
		if line, err := br.ReadString('\n'); err != nil {
			r.Close()
			return nil, err
		} else if m := re.FindStringSubmatch(line); m != nil {
			go io.Copy(ioutil.Discard, br)
			return m, nil
		}
	}
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
