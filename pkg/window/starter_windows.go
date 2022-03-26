//go:build windows
// +build windows

package window

import "syscall"

func getWindowWidth() int {

	return getSystemMetrics(0)
}
func getWindowHeight() int {

	return getSystemMetrics(1)
}

func getSystemMetrics(nIndex int) int {
	ret, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(uintptr(nIndex))
	return int(ret)
}
