package ssh

import (
	"bytes"
	"teamide/internal/context"
)

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

func (this_ *ShellClient) processZModem(buff []byte, n int, buffSize int) (isZModem bool, err error) {
	isZModem = true
	if this_.ZModemSZOO {
		this_.ZModemSZOO = false
		// 经过测试 centos7-8 使用的 lrzsz-0.12.20 在 sz 结束时会发送 ZModemSZEndOO
		// 而 deepin20 等自带更新的 lrzsz-0.12.21rc 在 sz 结束时不会发送 ZModemSZEndOO， 而前端 zmodemjs
		// 库只有接收到 ZModemSZEndOO 才会认为 sz 结束，固这里需判断 sz 结束时是否发送了 ZModemSZEndOO，
		// 如果没有则手动发送一个，以便保证前端 zmodemjs 库正常运行（如果不发送，会导致使用 sz 命令时无法连续
		// 下载多个文件）。
		if n < 2 {
			// 手动发送 ZModemSZEndOO
			this_.WSWriteBinary(ZModemSZEndOO)
			this_.WSWriteStdout(string(buff[:n]))
		} else if n == 2 {
			if buff[0] == ZModemSZEndOO[0] && buff[1] == ZModemSZEndOO[1] {
				this_.WSWriteBinary(ZModemSZEndOO)
			} else {
				// 手动发送 ZModemSZEndOO
				this_.WSWriteBinary(ZModemSZEndOO)
				this_.WSWriteStdout(string(buff[:n]))
			}
		} else {
			if buff[0] == ZModemSZEndOO[0] && buff[1] == ZModemSZEndOO[1] {
				this_.WSWriteBinary(buff[:2])
				this_.WSWriteStdout(string(buff[2:n]))
			} else {
				// 手动发送 ZModemSZEndOO
				this_.WSWriteBinary(ZModemSZEndOO)
				this_.WSWriteStdout(string(buff[:n]))
			}
		}
	} else {
		if this_.ZModemSZ {
			if (n) == buffSize {
				// 如果读取的长度为 buffsize，则认为是在传输数据，
				// 这样可以提高 sz 下载速率，很低概率会误判 zmodem 取消操作
				this_.WSWriteBinary(buff[:n])
			} else {
				if x, ok := ByteContains(buff[:n], ZModemSZEnd); ok {
					this_.ZModemSZ = false
					this_.ZModemSZOO = true
					this_.WSWriteBinary(ZModemSZEnd)
					if len(x) != 0 {
						this_.WSWriteConsole(string(x))
					}
				} else if _, ok := ByteContains(buff[:n], ZModemCancel); ok {
					this_.ZModemSZ = false
					this_.WSWriteBinary(buff[:n])
				} else {
					this_.WSWriteBinary(buff[:n])
				}
			}
		} else if this_.ZModemRZ {
			if x, ok := ByteContains(buff[:n], ZModemRZEnd); ok {
				out := map[string]interface{}{
					"token":    this_.Token,
					"fileSize": this_.rzFileSize,
					"isEnd":    true,
				}

				this_.rzFileSize = 0
				this_.rzFileUploadSize = 0

				this_.ZModemRZ = false
				this_.WSWriteBinary(ZModemRZEnd)
				if len(x) != 0 {
					this_.WSWriteConsole(string(x))
				}
				context.ServerWebsocketOutEvent("ssh-rz-upload", out)
			} else if _, ok := ByteContains(buff[:n], ZModemCancel); ok {
				out := map[string]interface{}{
					"token":    this_.Token,
					"fileSize": this_.rzFileSize,
					"isEnd":    true,
				}
				this_.rzFileSize = 0
				this_.rzFileUploadSize = 0

				this_.ZModemRZ = false
				this_.WSWriteBinary(buff[:n])

				context.ServerWebsocketOutEvent("ssh-rz-upload", out)
			} else {
				// rz 上传过程中服务器端还是会给客户端发送一些信息，比如心跳
				//this_.ZModemWriteJSON(&message{Type: messageTypeConsole, Data: buff[:n]})
				//this_.ZModemWriteMessage(websocket.BinaryMessage, buff[:n])

				startIndex := bytes.Index(buff[:n], ZModemRZCtrlStart)
				if startIndex != -1 {
					endIndex := bytes.Index(buff[:n], ZModemRZCtrlEnd1)
					if endIndex != -1 {
						ctrl := append(ZModemRZCtrlStart, buff[startIndex+len(ZModemRZCtrlStart):endIndex]...)
						ctrl = append(ctrl, ZModemRZCtrlEnd1...)
						this_.WSWriteBinary(ctrl)
						info := append(buff[:startIndex], buff[endIndex+len(ZModemRZCtrlEnd1):n]...)
						if len(info) != 0 {
							this_.WSWriteConsole(string(info))
						}
					} else {
						endIndex = bytes.Index(buff[:n], ZModemRZCtrlEnd2)
						if endIndex != -1 {
							ctrl := append(ZModemRZCtrlStart, buff[startIndex+len(ZModemRZCtrlStart):endIndex]...)
							ctrl = append(ctrl, ZModemRZCtrlEnd2...)
							this_.WSWriteBinary(ctrl)
							info := append(buff[:startIndex], buff[endIndex+len(ZModemRZCtrlEnd2):n]...)
							if len(info) != 0 {
								this_.WSWriteConsole(string(info))
							}
						} else {
							this_.WSWriteConsole(string(buff[:n]))
						}
					}
				} else {
					this_.WSWriteConsole(string(buff[:n]))
				}
			}
		} else {
			if x, ok := ByteContains(buff[:n], ZModemSZStart); ok {
				if this_.DisableZModemSZ {
					this_.WSWriteAlert("sz download is disabled")
					this_.ZModemWriteSSH(ZModemCancel)
				} else {
					if y, ok := ByteContains(x, ZModemCancel); ok {
						// 下载不存在的文件以及文件夹(zmodem 不支持下载文件夹)时
						this_.WSWriteStdout(string(y))
					} else {
						this_.ZModemSZ = true
						if len(x) != 0 {
							this_.WSWriteConsole(string(x))
						}
						this_.WSWriteBinary(ZModemSZStart)
					}
				}
			} else if x, ok := ByteContains(buff[:n], ZModemRZStart); ok {
				if this_.DisableZModemRZ {
					this_.WSWriteAlert("rz upload is disabled")
					this_.ZModemWriteSSH(ZModemCancel)
				} else {
					this_.ZModemRZ = true
					this_.WSWriteEvent("shell to upload file")
					if len(x) != 0 {
						this_.WSWriteConsole(string(x))
					}
					this_.WSWriteBinary(ZModemRZStart)
				}
			} else if x, ok := ByteContains(buff[:n], ZModemRZEStart); ok {
				if this_.DisableZModemRZ {
					this_.WSWriteAlert("rz upload is disabled")
					this_.ZModemWriteSSH(ZModemCancel)
				} else {
					this_.ZModemRZ = true
					if len(x) != 0 {
						this_.WSWriteConsole(string(x))
					}
					this_.WSWriteBinary(ZModemRZEStart)
				}
			} else if x, ok := ByteContains(buff[:n], ZModemRZSStart); ok {
				if this_.DisableZModemRZ {
					this_.WSWriteAlert("rz upload is disabled")
					this_.ZModemWriteSSH(ZModemCancel)
				} else {
					this_.ZModemRZ = true
					if len(x) != 0 {
						this_.WSWriteConsole(string(x))
					}
					this_.WSWriteBinary(ZModemRZSStart)
				}
			} else if x, ok := ByteContains(buff[:n], ZModemRZESStart); ok {
				if this_.DisableZModemRZ {
					this_.ZModemWriteSSH(ZModemCancel)
				} else {
					this_.ZModemRZ = true
					if len(x) != 0 {
						this_.WSWriteConsole(string(x))
					}
					this_.WSWriteBinary(ZModemRZESStart)
				}
			} else {
				isZModem = false
			}
		}
	}
	return
}

func (this_ *ShellClient) ZModemWriteSSH(message []byte) {
	this_.SSHWrite(message)
}

func ByteContains(x, y []byte) (n []byte, contain bool) {
	index := bytes.Index(x, y)
	if index == -1 {
		return
	}
	lastIndex := index + len(y)
	n = append(x[:index], x[lastIndex:]...)
	return n, true
}
