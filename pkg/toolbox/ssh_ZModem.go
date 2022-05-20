package toolbox

var (
	//ZModemSZStart sz fmt.Sprintf("%+q", "rz\r**\x18B00000000000000\r\x8a\x11")
	//ZModemSZStart = []byte{13, 42, 42, 24, 66, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 13, 138, 17}
	ZModemSZStart = []byte{42, 42, 24, 66, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 13, 138, 17}
	//ZModemSZEnd sz 结束 fmt.Sprintf("%+q", "\r**\x18B0800000000022d\r\x8a")
	//ZModemSZEnd = []byte{13, 42, 42, 24, 66, 48, 56, 48, 48, 48, 48, 48, 48, 48, 48, 48, 50, 50, 100, 13, 138}
	ZModemSZEnd = []byte{42, 42, 24, 66, 48, 56, 48, 48, 48, 48, 48, 48, 48, 48, 48, 50, 50, 100, 13, 138}
	//ZModemSZEndOO sz 结束后可能还会发送两个 OO，但是经过测试发现不一定每次都会发送 fmt.Sprintf("%+q", "OO")
	ZModemSZEndOO = []byte{79, 79}

	//ZModemRZStart rz fmt.Sprintf("%+q", "**\x18B0100000023be50\r\x8a\x11")
	ZModemRZStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 48, 50, 51, 98, 101, 53, 48, 13, 138, 17}
	//ZModemRZEStart rz -e fmt.Sprintf("%+q", "**\x18B0100000063f694\r\x8a\x11")
	ZModemRZEStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 48, 54, 51, 102, 54, 57, 52, 13, 138, 17}
	//ZModemRZSStart rz -S fmt.Sprintf("%+q", "**\x18B0100000223d832\r\x8a\x11")
	ZModemRZSStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 50, 50, 51, 100, 56, 51, 50, 13, 138, 17}
	//ZModemRZESStart rz -e -S fmt.Sprintf("%+q", "**\x18B010000026390f6\r\x8a\x11")
	ZModemRZESStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 50, 54, 51, 57, 48, 102, 54, 13, 138, 17}
	//ZModemRZEnd rz 结束 fmt.Sprintf("%+q", "**\x18B0800000000022d\r\x8a")
	ZModemRZEnd = []byte{42, 42, 24, 66, 48, 56, 48, 48, 48, 48, 48, 48, 48, 48, 48, 50, 50, 100, 13, 138}

	//ZModemRZCtrlStart **\x18B0
	ZModemRZCtrlStart = []byte{42, 42, 24, 66, 48}
	//ZModemRZCtrlEnd1 \r\x8a\x11
	ZModemRZCtrlEnd1 = []byte{13, 138, 17}
	//ZModemRZCtrlEnd2 \r\x8a
	ZModemRZCtrlEnd2 = []byte{13, 138}

	//ZModemCancel zmodem 取消 \x18\x18\x18\x18\x18\x08\x08\x08\x08\x08
	ZModemCancel = []byte{24, 24, 24, 24, 24, 8, 8, 8, 8, 8}
)

type ZModemMessageType string

const (
	ZModemMessageTypeAddr      ZModemMessageType = "addr"
	ZModemMessageTypeTerm      ZModemMessageType = "term"
	ZModemMessageTypeLogin     ZModemMessageType = "login"
	ZModemMessageTypePassword  ZModemMessageType = "password"
	ZModemMessageTypePublickey ZModemMessageType = "publickey"
	ZModemMessageTypeStdin     ZModemMessageType = "stdin"
	ZModemMessageTypeStdout    ZModemMessageType = "stdout"
	ZModemMessageTypeStderr    ZModemMessageType = "stderr"
	ZModemMessageTypeResize    ZModemMessageType = "resize"
	ZModemMessageTypeIgnore    ZModemMessageType = "ignore"
	ZModemMessageTypeConsole   ZModemMessageType = "console"
	ZModemMessageTypeAlert     ZModemMessageType = "alert"
)

type ZModemMessage struct {
	Type ZModemMessageType `json:"type"`
	Data []byte            `json:"data,omitempty"`
	Cols int               `json:"cols,omitempty"`
	Rows int               `json:"rows,omitempty"`
}
