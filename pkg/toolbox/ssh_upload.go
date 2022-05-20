package toolbox

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
)

func (this_ *SSHShellClient) processZModem(buff []byte, n int, buffSize int) (isZModem bool, err error) {
	isZModem = true
	t := ZModemMessageTypeStdout
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
			this_.ZModemWriteJSON(&ZModemMessage{Type: t, Data: buff[:n]})
		} else if n == 2 {
			if buff[0] == ZModemSZEndOO[0] && buff[1] == ZModemSZEndOO[1] {
				this_.WSWriteBinary(ZModemSZEndOO)
			} else {
				// 手动发送 ZModemSZEndOO
				this_.WSWriteBinary(ZModemSZEndOO)
				this_.ZModemWriteJSON(&ZModemMessage{Type: t, Data: buff[:n]})
			}
		} else {
			if buff[0] == ZModemSZEndOO[0] && buff[1] == ZModemSZEndOO[1] {
				this_.WSWriteBinary(buff[:2])
				this_.ZModemWriteJSON(&ZModemMessage{Type: t, Data: buff[2:n]})
			} else {
				// 手动发送 ZModemSZEndOO
				this_.WSWriteBinary(ZModemSZEndOO)
				this_.ZModemWriteJSON(&ZModemMessage{Type: t, Data: buff[:n]})
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
						this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeConsole, Data: x})
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
				this_.ZModemRZ = false
				this_.WSWriteBinary(ZModemRZEnd)
				if len(x) != 0 {
					this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeConsole, Data: x})
				}
			} else if _, ok := ByteContains(buff[:n], ZModemCancel); ok {
				this_.ZModemRZ = false
				this_.WSWriteBinary(buff[:n])
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
							this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeConsole, Data: info})
						}
					} else {
						endIndex = bytes.Index(buff[:n], ZModemRZCtrlEnd2)
						if endIndex != -1 {
							ctrl := append(ZModemRZCtrlStart, buff[startIndex+len(ZModemRZCtrlStart):endIndex]...)
							ctrl = append(ctrl, ZModemRZCtrlEnd2...)
							this_.WSWriteBinary(ctrl)
							info := append(buff[:startIndex], buff[endIndex+len(ZModemRZCtrlEnd2):n]...)
							if len(info) != 0 {
								this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeConsole, Data: info})
							}
						} else {
							this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeConsole, Data: buff[:n]})
						}
					}
				} else {
					this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeConsole, Data: buff[:n]})
				}
			}
		} else {
			if x, ok := ByteContains(buff[:n], ZModemSZStart); ok {
				if this_.DisableZModemSZ {
					this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeAlert, Data: []byte("sz download is disabled")})
					this_.ZModemWriteSSH(ZModemCancel)
				} else {
					if y, ok := ByteContains(x, ZModemCancel); ok {
						// 下载不存在的文件以及文件夹(zmodem 不支持下载文件夹)时
						this_.ZModemWriteJSON(&ZModemMessage{Type: t, Data: y})
					} else {
						this_.ZModemSZ = true
						if len(x) != 0 {
							this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeConsole, Data: x})
						}
						this_.WSWriteBinary(ZModemSZStart)
					}
				}
			} else if x, ok := ByteContains(buff[:n], ZModemRZStart); ok {
				if this_.DisableZModemRZ {
					this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeAlert, Data: []byte("rz upload is disabled")})
					this_.ZModemWriteSSH(ZModemCancel)
				} else {
					this_.ZModemRZ = true
					this_.WSWriteEvent("shell to upload file")
					if len(x) != 0 {
						this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeConsole, Data: x})
					}
					this_.WSWriteBinary(ZModemRZStart)
				}
			} else if x, ok := ByteContains(buff[:n], ZModemRZEStart); ok {
				if this_.DisableZModemRZ {
					this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeAlert, Data: []byte("rz upload is disabled")})
					this_.ZModemWriteSSH(ZModemCancel)
				} else {
					this_.ZModemRZ = true
					if len(x) != 0 {
						this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeConsole, Data: x})
					}
					this_.WSWriteBinary(ZModemRZEStart)
				}
			} else if x, ok := ByteContains(buff[:n], ZModemRZSStart); ok {
				if this_.DisableZModemRZ {
					this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeAlert, Data: []byte("rz upload is disabled")})
					this_.ZModemWriteSSH(ZModemCancel)
				} else {
					this_.ZModemRZ = true
					if len(x) != 0 {
						this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeConsole, Data: x})
					}
					this_.WSWriteBinary(ZModemRZSStart)
				}
			} else if x, ok := ByteContains(buff[:n], ZModemRZESStart); ok {
				if this_.DisableZModemRZ {
					this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeAlert, Data: []byte("rz upload is disabled")})
					this_.ZModemWriteSSH(ZModemCancel)
				} else {
					this_.ZModemRZ = true
					if len(x) != 0 {
						this_.ZModemWriteJSON(&ZModemMessage{Type: ZModemMessageTypeConsole, Data: x})
					}
					this_.WSWriteBinary(ZModemRZESStart)
				}
			} else {
				isZModem = false
				//this_.ZModemWriteJSON(&ZModemMessage{Type: t, Data: buff[:n]})
			}
		}
	}
	return
}

func (this_ *SSHShellClient) ZModemWriteSSH(message []byte) {
	this_.SSHWrite(message)
}
func (this_ *SSHShellClient) ZModemWriteJSON(message *ZModemMessage) {

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

var (

	// RZStart rz
	RZStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 48, 50, 51, 98, 101, 53, 48, 13, 138, 17}
	// RZEStart rz -e
	RZEStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 48, 54, 51, 102, 54, 57, 52, 13, 138, 17}
)

func (this_ *SSHShellClient) listenUpload() {
	if this_.UploadFile == nil {
		this_.UploadFile = make(chan *UploadFile, 10)

		go func() {
			for {
				select {
				case uploadFile := <-this_.UploadFile:
					this_.upload(uploadFile)
				}
			}

		}()
	}
	return
}

func (this_ *SSHShellClient) upload(uploadFile *UploadFile) {

	return
}

func SSHUpload(c *gin.Context) (res interface{}, err error) {
	token := c.PostForm("token")
	if token == "" {
		err = errors.New("token获取失败")
		return
	}
	client := SSHShellCache[token]
	if client == nil {
		err = errors.New("SSH会话丢失")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		return
	}

	uploadFile := &UploadFile{
		File:     file,
		FullPath: c.PostForm("fullPath"),
	}
	client.UploadFile <- uploadFile

	return
}

func SSHDownload(data map[string]string, c *gin.Context) (err error) {

	token := data["token"]
	if token == "" {
		err = errors.New("token获取失败")
		return
	}
	place := data["place"]
	if place == "" {
		err = errors.New("place获取失败")
		return
	}
	path := data["path"]
	if path == "" {
		err = errors.New("path获取失败")
		return
	}
	client := SSHSftpCache[token]
	if client == nil {
		err = errors.New("SSH会话丢失")
		return
	}
	if place == "local" {
		err = client.localDownload(c, path)
	} else if place == "remote" {
		err = client.remoteDownload(c, path)
	}

	return
}
